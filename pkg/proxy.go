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
)

var ErrDoesNotExist = errors.New("item does not exist")

type ProxGe struct {
	router *mux.Router
	cache  GECache
	apis   []GEApi
}

func New(cache GECache, router *mux.Router, api ...GEApi) *ProxGe {
	p := &ProxGe{
		router: router,
		cache:  cache,
		apis:   api,
	}

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/id/{id}", asJson(p.getById))
	return p
}

type priceResponse struct {
	ID    int `json:"id"`
	Price int `json:"price"`
}

func (p *ProxGe) getById(w http.ResponseWriter, r *http.Request) {
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

	price, err := p.cache.Get(id)
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
