package transformer

import (
	"sync"
)

type Transformer interface {
	Transform(wg *sync.WaitGroup, jobs chan []byte, results chan int)
}
