package main

import (
	"sync"
)

type (

	// ChannelInterface is a channel which accepts all kind of data
	ChannelInterface chan interface{}

	// Dispatcher has an input channel and send data to other data
	Dispatcher struct {
		sync.RWMutex
		InputChannel   ChannelInterface
		OutputChannels []ChannelInterface
	}
)

var GlobalDispatcher *Dispatcher

func init() {
	GlobalDispatcher = NewDispatcher()
}

// NewDispatcher returns a new Dispatcher with only 10 output channels
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		InputChannel:   make(ChannelInterface),
		OutputChannels: make([]ChannelInterface, 10),
	}
}

// Run waits on InputChannel and send data to OutputChannels
func (d *Dispatcher) Run() {

	var i interface{}

	for {

		i = <-d.InputChannel
		d.RLock()
		for _, channel := range d.OutputChannels {
			if channel != nil {
				channel <- i
			}
		}
		d.RUnlock()

	}

}

// GetChannel returns a channel where data will be sent, if no channel is available, nil will be returned
func (d *Dispatcher) GetChannel() ChannelInterface {
	d.Lock()
	defer d.Unlock()

	for i, channel := range d.OutputChannels {
		if channel != nil {
			d.OutputChannels[i] = make(ChannelInterface)
			return d.OutputChannels[i]
		}
	}

	return nil

}
