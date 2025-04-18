package exchange

import (
	"fmt"
	"sync"
)

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
	err error
}

func (ex *Exchange) GetError() error {
	return ex.err
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
	if _, ok := ex.Exc[idChannel]; !ok {
		ex.Exc[idChannel] = make(chan StatusChannel, defaultSize)
	} else {
		ex.err = fmt.Errorf("failed create channel")
	}
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
	go func() {
		ex.mu.Lock()
		defer ex.mu.Unlock()
		ex.Exc[idChannel] <- status
	}()
}
