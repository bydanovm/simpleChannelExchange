package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bydanovm/simpleChannelExchange/internal/exchange"
)

func main() {
	// Initializing two channels
	exc := exchange.Init()
	// Init channel exchange with buffer size 5
	if err := exc.NewChannel(exchange.ExchangeChannel1, 5).GetError(); err != nil {
		panic(fmt.Errorf("%w", err).Error())
	}
	// Init channel exchange with default buffer size 10
	if err := exc.NewChannel(exchange.ExchangeChannel2).GetError(); err != nil {
		panic(fmt.Errorf("%w", err).Error())
	}

	// Writing a message to channel 1
	go func() {
		for i := 0; i < 10; i++ {
			exc.WriteChannel(
				exchange.ExchangeChannel1,
				exchange.StatusChannel{
					Data: fmt.Sprintf("message id:%d for channel 1", i),
				},
			)
		}
	}()

	// Writing a message to channel 2
	go func() {
		for i := 0; i < 10; i++ {
			exc.WriteChannel(
				exchange.ExchangeChannel2,
				exchange.StatusChannel{
					Data: fmt.Sprintf("message id:%d for channel 2", i),
				},
			)
		}
	}()

	// Reading a message into channel 1
	go func() {
		for rcvMSg := range exc.ReadChannel(exchange.ExchangeChannel1) {
			if err := rcvMSg.GetError(); err != nil {
				fmt.Println(fmt.Errorf("error revieve message: %w", err).Error())
			} else {
				if v, ok := rcvMSg.Data.(string); ok {
					fmt.Printf("Recieve message from channel 1: %s\n", v)
				} else {
					fmt.Printf("error type assertion for message: %v\n", rcvMSg.Data)
				}
			}
		}
	}()

	// Reading a message into channel 1
	go func() {
		for rcvMSg := range exc.ReadChannel(exchange.ExchangeChannel2) {
			if v, ok := rcvMSg.Data.(string); ok {
				fmt.Printf("Recieve message from channel 2: %s\n", v)
			} else {
				fmt.Printf("Error revieve message: %v\n", rcvMSg.Data)
			}
		}
	}()

	for {
		<-time.After(time.Second * 5)
		os.Exit(1)
	}
}
