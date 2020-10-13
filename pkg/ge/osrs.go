package ge

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/while-loop/proxge/pkg"
	"net/http"
	"time"
)

type GEResp struct {
	Item struct {
		Icon        string `json:"icon"`
		IconLarge   string `json:"icon_large"`
		ID          int    `json:"id"`
		Type        string `json:"type"`
		TypeIcon    string `json:"typeIcon"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Current     struct {
			Trend string `json:"trend"`
			Price int    `json:"price"`
		} `json:"current"`
		Today struct {
			Trend string `json:"trend"`
			Price string `json:"price"`
		} `json:"today"`
		Members string `json:"members"`
		Day30   struct {
			Trend  string `json:"trend"`
			Change string `json:"change"`
		} `json:"day30"`
		Day90 struct {
			Trend  string `json:"trend"`
			Change string `json:"change"`
		} `json:"day90"`
		Day180 struct {
			Trend  string `json:"trend"`
			Change string `json:"change"`
		} `json:"day180"`
	} `json:"item"`
}

const baseUrl = "https://secure.runescape.com/m=itemdb_oldschool/api/catalogue/detail.json?item=%d"

type osrsGe struct {
	client *http.Client
}

func NewOsrsGe() proxge.GEApi {
	return &osrsGe{
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (ge *osrsGe) PriceById(id int) (int, error) {
	resp, err := ge.client.Get(fmt.Sprintf(baseUrl, id))
	if err != nil {
		return 0, errors.Wrap(err, "osrs")
	}

	defer resp.Body.Close()

	var data *GEResp
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return 0, errors.Wrap(err, "osrs")
	}

	if data.Item.Current.Price <= 0 {
		return 0, fmt.Errorf("no price given for item %d %s", data.Item.ID, data.Item.Name)
	}

	return data.Item.Current.Price, nil
}

func (ge *osrsGe) Name() string {
	return "osrsge"
}
