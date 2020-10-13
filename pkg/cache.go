package proxge

import "time"

type GECache interface {
	Get(id int) (int, time.Duration, error)
	Set(id int, price int) error
	GetTTL() time.Duration
}
