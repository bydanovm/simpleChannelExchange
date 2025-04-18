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

type StatusChannelImpl interface {
	GetError() error
	GetData() interface{}
}
type StatusChannel struct {
	Data interface{}
	err  error
}

func (sc StatusChannel) GetError() error {
	return sc.err
}

func (sc StatusChannel) GetData() interface{} {
	return sc.Data
}

type Exchange struct {
	mu  sync.RWMutex
	Exc map[interface{}]chan StatusChannelImpl
	err error
}

func (ex *Exchange) GetError() error {
	return ex.err
}

func Init() *Exchange {
	exc := &Exchange{
		Exc: make(map[interface{}]chan StatusChannelImpl),
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
		ex.Exc[idChannel] = make(chan StatusChannelImpl, defaultSize)
	} else {
		ex.err = fmt.Errorf("failed create channel")
	}
	return ex
}

func (ex *Exchange) ReadChannel(idChannel int) <-chan StatusChannelImpl {
	var outCh = make(chan StatusChannelImpl, cap(ex.Exc[idChannel]))
	go func() {
		for s := range ex.Exc[idChannel] {
			outCh <- s
		}
	}()
	return outCh
}

func (ex *Exchange) WriteChannel(idChannel int, status StatusChannelImpl) {
	go func() {
		ex.mu.Lock()
		defer ex.mu.Unlock()
		ex.Exc[idChannel] <- status
	}()
}
