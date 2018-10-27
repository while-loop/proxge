package proxge

type GECache interface {
	Get(id int) (int, error)
	Set(id int, price int) error
}
