package main

import (
	"fmt"
	"regexp"
	"sync"
	"time"
)

const (
	LenChannels = 100
)

var GlobalCollectors []Collector

func init() {
	GlobalCollectors = []Collector{
		NewCollectDateTimeRequest(),
		NewCollectUrl(),
	}
}

type (
	// Date is a candy struct for string
	Date string

	// HoursHits is a candy struct to keep hits for every hour
	HoursHits map[string]int

	// CollectDateTimeRequests collects for every day how many hits per hours
	CollectDateTimeRequests struct {
		Input     chan *LogLineStruct
		CountData map[Date]HoursHits
	}

	// CollectUrl collects urls and the most visited urls
	CollectUrl struct {
		Input     chan *LogLineStruct
		CountData map[string]int
	}

	// Collector defines a generic collector
	// A collector has a run function to receive the data
	// end expone a channel where data will be sent
	Collector interface {
		Run(<-chan time.Time, *sync.WaitGroup)
		GetChannel() chan *LogLineStruct
	}
)

// NewCollectDateTimeRequest initializes the struct
func NewCollectDateTimeRequest() *CollectDateTimeRequests {
	return &CollectDateTimeRequests{
		Input:     make(chan *LogLineStruct, LenChannels),
		CountData: make(map[Date]HoursHits),
	}
}

// Run runs the an infinity loop for make things with data
// Every tick prints the current status of data
func (c *CollectDateTimeRequests) Run(tick <-chan time.Time, wait *sync.WaitGroup) {

	defer wait.Done()

	var line *LogLineStruct
	var ok bool

	for {

		select {

		case line, ok = <-c.Input:
			if !ok {
				return
			}

			date := Date(line.Date.Format("2006-01-02")) // cast to Date format
			if _, ok := c.CountData[date]; !ok {
				c.CountData[date] = make(HoursHits)
			}
			c.CountData[date][line.Date.Format("15")]++

		case <-tick:
			fmt.Println(c.CountData)
		}

	}

}

// GetChannel returns the channel where send data
func (c *CollectDateTimeRequests) GetChannel() chan *LogLineStruct {
	return c.Input
}

func NewCollectUrl() *CollectUrl {
	return &CollectUrl{
		Input:     make(chan *LogLineStruct, LenChannels),
		CountData: make(map[string]int),
	}
}

// Run runs the an infinity loop for make things with data
// Every tick prints the current status of data
func (c *CollectUrl) Run(tick <-chan time.Time, wait *sync.WaitGroup) {

	defer wait.Done()

	var line *LogLineStruct
	var ok bool

	re := regexp.MustCompile("(wp-|mp4|feed)")

	for {

		select {

		case line, ok = <-c.Input:
			if !ok {
				return
			}

			if line.Method == "GET" && !re.MatchString(line.URL) {
				c.CountData[line.URL]++
			}

		case <-tick:
			fmt.Println(TopHits(c.CountData, 10))
		}

	}

}

// GetChannel returns the channel where send data
func (c *CollectUrl) GetChannel() chan *LogLineStruct {
	return c.Input
}
