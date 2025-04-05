package filter

const (
	numWorkers = 100
)

type Filter interface {
	Execute()
}
