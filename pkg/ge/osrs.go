package ge

import (
	"github.com/while-loop/proxge/pkg"
	"net/http"
	"time"
)

type osrsGe struct {
	client *http.Client
}

func NewOsrsGe() proxge.GEApi {
	return &osrsGe{
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (ge *osrsGe) PriceById(id int) (int, error) {
	panic("implement me")
}
