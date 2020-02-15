package ge

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/while-loop/proxge/pkg"
	"log"
	"net/http"
	"time"
)

const summaryUrl = `https://rsbuddy.com/exchange/summary.json`

var _ proxge.GEApi = &rsBuddyGe{}
type rsBuddyGe struct {
	client    *http.Client
	cache     proxge.GECache
	lastCache time.Time
}

type item struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Members         bool   `json:"members"`
	Sp              int    `json:"sp"`
	BuyAverage      int    `json:"buy_average"`
	BuyQuantity     int    `json:"buy_quantity"`
	SellAverage     int    `json:"sell_average"`
	SellQuantity    int    `json:"sell_quantity"`
	OverallAverage  int    `json:"overall_average"`
	OverallQuantity int    `json:"overall_quantity"`
}

func NewRsBuddyGe(cache proxge.GECache) proxge.GEApi {
	return &rsBuddyGe{
		client:    &http.Client{Timeout: 5 * time.Second},
		cache:     cache,
		lastCache: time.Unix(0, 0),
	}
}

func (ge *rsBuddyGe) PriceById(id int) (int, error) {
	resp, err := ge.client.Get(summaryUrl)
	if err != nil {
		return 0, errors.Wrap(err, "osbuddy")
	}

	defer resp.Body.Close()

	prices := map[string]item{}
	if err := json.NewDecoder(resp.Body).Decode(&prices); err != nil {
		return 0, errors.Wrap(err, "osbuddy")
	}

	if ge.lastCache.Add(10 * time.Minute).Before(time.Now()) {
		ge.lastCache = time.Now()
		go func() {
			i := 0
			for _, v := range prices {
				log.Println(v)
				i++
				if err := ge.cache.Set(v.ID, v.OverallAverage); err != nil {
					log.Printf("osbuddy: failed to set price for %d: %v", v.ID, err)
				}
				if i > 50 {
					break
				}
			}
		}()
	}

	item, ok := prices[fmt.Sprintf("%d", id)]
	if !ok {
		return 0, proxge.ErrDoesNotExist
	}

	return item.OverallAverage, nil
}

func (ge *rsBuddyGe) Name() string {
	return "rsbuddyge"
}
