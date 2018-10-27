package proxge

type GEApi interface {
	PriceById(id int) (int, error)
	Name() string
}
