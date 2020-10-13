package proxge

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var ErrDoesNotExist = errors.New("item does not exist")

type ProxGe struct {
	cache  GECache
	apis   []GEApi
}

func New(cache GECache, router *mux.Router, api ...GEApi) *ProxGe {
	p := &ProxGe{
		cache:  cache,
		apis:   api,
	}

	router.HandleFunc("/id/{id}", asJson(p.GetById))
	return p
}

type priceResponse struct {
	ID    int `json:"id"`
	Price int `json:"price"`
}

func (p *ProxGe) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strId, ok := vars["id"]
	if !ok {
		writeError(http.StatusBadRequest, "no id", w)
		return
	}
	strId = strings.TrimSpace(strId)

	id, err := strconv.Atoi(strId)
	if err != nil {
		writeError(http.StatusBadRequest, "bad int: "+err.Error(), w)
		return
	}

	price, ttl, err := p.cache.Get(id)
	if err != nil {
		price, err = p.getPrice(id)
		if err != nil {
			writeError(http.StatusInternalServerError, err.Error(), w)
			return
		}

		err := p.cache.Set(id, price)
		if err != nil {
			log.Printf("failed to set cache for %d %d: %v", id, price, err)
		}
	}
	now := time.Now().UTC()
	cacheSince := now.Add(ttl-p.cache.GetTTL()).Format(http.TimeFormat)
	cacheUntil := now.Add(ttl).Format(http.TimeFormat)

	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, public", int(ttl.Seconds())))
	w.Header().Set("Last-Modified", cacheSince)
	w.Header().Set("Expires", cacheUntil)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(priceResponse{Price: price, ID: id}); err != nil {
		log.Printf("failed to send price response for %d: %v\n", id, err)
	}
}

func (p *ProxGe) getPrice(id int) (int, error) {
	for _, api := range p.apis {
		price, err := api.PriceById(id)
		if err == nil {
			return price, nil
		}

		log.Printf("failed to get price for %d using %T: %v\n", id, api, err)
	}

	return 0, fmt.Errorf("unable to get price for item: %d", id)
}
