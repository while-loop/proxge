package ge

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/while-loop/proxge/pkg"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type GEResp struct {
	Item struct {
		ID      int                    `json:"id"`
		Name    string                 `json:"name"`
		Current map[string]interface{} `json:"current"`
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
		return -1, errors.Wrap(err, "osrs")
	}

	defer resp.Body.Close()

	var data *GEResp
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return -1, errors.Wrap(err, "osrs")
	}

	priceStr, ok := data.Item.Current["price"]
	if !ok || priceStr == "" {
		return -1, fmt.Errorf("no price given for item %d %s", data.Item.ID, data.Item.Name)
	}

	price := 0
	switch priceStr.(type) {
	case float64:
		price = int(priceStr.(float64))
	case string:
		price = unhumanizeNumber(priceStr.(string))
	default:
	}

	if price <= 0 {
		return -1, fmt.Errorf("no price given for item %d %s", data.Item.ID, data.Item.Name)
	}

	return price, nil
}

func (ge *osrsGe) Name() string {
	return "osrsge"
}

func unhumanizeNumber(num string) int {
	num = strings.ReplaceAll(strings.ToLower(num), ",", "")
	factor := 1.0
	if strings.Contains(num, "k") {
		factor = 1000
	} else if strings.Contains(num, "m") {
		factor = 1000000
	} else if strings.Contains(num, "b") {
		factor = 1000000000
	}

	num = strings.TrimRight(strings.TrimSpace(num), "kmb")
	numFloat, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return 0
	}

	numFloat = numFloat * factor
	numFloat = math.Round(numFloat)

	return int(numFloat)
}
