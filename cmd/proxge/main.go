package main

import (
	"flag"
	"github.com/akrylysov/algnhsa"
	"github.com/go-redis/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/while-loop/proxge/pkg"
	"github.com/while-loop/proxge/pkg/cache"
	"github.com/while-loop/proxge/pkg/ge"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	RedisAddr string `split_words:"true"`
	RedisPass string `split_words:"true"`
	RedisDB   int    `split_words:"true" default:"0"`
	Laddr     string
	Exp       time.Duration `default:"30m"`
}

var (
	v = flag.Bool("v", false, proxge.Name+" version")
)

func main() {
	log.Printf("%s %s %s %s", proxge.Name, proxge.Version, proxge.BuildTime, proxge.Commit)

	flag.Parse()
	if *v {
		return
	}

	var config Config
	err := envconfig.Process("proxge", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	c := cache.NewMemCache()
	if config.RedisAddr != "" {
		c = cache.NewRedisCache(cache.NewRedisClient(&redis.Options{
			Addr:     config.RedisAddr,
			DB:       config.RedisDB,
			Password: config.RedisPass,
		}), config.Exp)
	}

	apis := []proxge.GEApi{
		ge.NewRsBuddyGe(c),
	}

	log.Printf("Using %T as cache\n", c)
	log.Printf("Using %v apis\n", c)

	router := mux.NewRouter()

	p := proxge.New(c, router, apis...)
	handler := wrapAppHandler(router)
	_ = p

	if config.Laddr != "" {
		log.Println("using local addr", config.Laddr)
		log.Println(http.ListenAndServe(config.Laddr, handler))
	} else {
		log.Println("Starting lambda")
		algnhsa.ListenAndServe(handler, nil)
	}
}

func wrapAppHandler(handler http.Handler) http.Handler {
	h := handlers.LoggingHandler(os.Stdout, handler)
	h = handlers.CORS()(h)
	h = handlers.RecoveryHandler()(h)
	return h
}
