package main

import (
	"regexp"
	"time"
)

const (
	LenChannels = 100
)

type (
	// Month is a candy struct for string
	Month string

	// Url is candy struct for string
	Url string

	// HoursHits is a candy struct to keep hits for every hour
	HoursHits map[string]int

	// CollectDateTimeRequests collects for every day how many hits per hours
	CollectDateTimeRequests struct {
		Input     chan LogLineStruct
		CountData map[Month]HoursHits
	}

	// CollectUrl collects urls and the most visited urls
	CollectUrl struct {
		Input     chan LogLineStruct
		CountData map[Url]int
	}

	// Collector defines a generic collector
	// A collector has a run function to receive the data
	// end expone a channel where data will be sent
	Collector interface {
		Run(<-chan time.Time)
		GetChannel() chan LogLineStruct
	}
)

var (
	// skyp generic pages in wordpress
	ExcludeRegex = "(wp-|mp4|feed|xmlrpc|all|css|sitemap.xml|downloads|" + regexp.QuoteMeta("/?s=") + ")"

	GlobalCollectors []Collector
)

func init() {

	GlobalCollectors = []Collector{
		NewCollectDateTimeRequest(),
		NewCollectUrl(),
	}

}

// NewCollectDateTimeRequest initializes the struct
func NewCollectDateTimeRequest() *CollectDateTimeRequests {
	return &CollectDateTimeRequests{
		Input:     make(chan LogLineStruct, LenChannels),
		CountData: make(map[Month]HoursHits),
	}
}

// Run runs the an infinity loop for make things with data
func (c *CollectDateTimeRequests) Run(tick <-chan time.Time) {

	var line LogLineStruct
	var ok bool

	re := regexp.MustCompile(ExcludeRegex)

	for {

		line, ok = <-c.Input
		if !ok {
			return
		}

		if line.Method == "GET" && !re.MatchString(line.URL) {
			date := Month(line.Date.Format("01")) // cast to Month format
			if _, ok := c.CountData[date]; !ok {
				c.CountData[date] = make(HoursHits)
			}
			c.CountData[date][line.Date.Format("15")]++
		}

	}

}

// GetChannel returns the channel where send data
func (c *CollectDateTimeRequests) GetChannel() chan LogLineStruct {
	return c.Input
}

// NewCollectUrl initializes the struct
func NewCollectUrl() *CollectUrl {
	return &CollectUrl{
		Input:     make(chan LogLineStruct, LenChannels),
		CountData: make(map[Url]int),
	}
}

// Run runs the an infinity loop for make things with data
func (c *CollectUrl) Run(tick <-chan time.Time) {

	var line LogLineStruct
	var ok bool

	re := regexp.MustCompile(ExcludeRegex)

	for {

		line, ok = <-c.Input
		if !ok {
			return
		}

		if line.Method == "GET" && line.URL != "/" && !re.MatchString(line.URL) {
			c.CountData[Url(line.URL)]++
		}

	}

}

// GetChannel returns the channel where send data
func (c *CollectUrl) GetChannel() chan LogLineStruct {
	return c.Input
}
