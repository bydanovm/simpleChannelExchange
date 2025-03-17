package exchange

import "sync"

const (
	ExchangeChannel1 int = iota + 1
	ExchangeChannel2
	ExchangeChannel3
	ExchangeChannel4
	ExchangeChannel5
)

var defaultSize int = 10

type Exchange struct {
	mu  sync.RWMutex
	Exc map[interface{}]chan StatusChannel
}

type StatusChannel struct {
	Data interface{}
	err  error
}

func (sc *StatusChannel) GetError() error {
	return sc.err
}

func Init() *Exchange {
	exc := &Exchange{
		Exc: make(map[interface{}]chan StatusChannel),
	}
	return exc
}

func (ex *Exchange) NewChannel(idChannel int, sizeChannel ...int) *Exchange {
	for size := range sizeChannel {
		defaultSize = size
		break
	}
	ex.mu.Lock()
	defer ex.mu.Unlock()
	ex.Exc[idChannel] = createChannel(defaultSize)
	return ex
}

func (ex *Exchange) ReadChannel(idChannel int) <-chan StatusChannel {
	var outCh = make(chan StatusChannel, cap(ex.Exc[idChannel]))
	go func() {
		for s := range ex.Exc[idChannel] {
			outCh <- s
		}
	}()
	return outCh
}

func (ex *Exchange) WriteChannel(idChannel int, status StatusChannel) {
	ex.mu.Lock()
	defer ex.mu.Unlock()
	ex.Exc[idChannel] <- status
}

func createChannel(sizeChannel ...int) chan StatusChannel {
	var defaultSize int = 10
	for size := range sizeChannel {
		defaultSize = size
		break
	}
	ch := make(chan StatusChannel, defaultSize)
	return ch
}
