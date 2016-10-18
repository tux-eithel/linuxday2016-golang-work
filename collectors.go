package main

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
)

// NewCollectDateTimeRequest initializes the struct
func NewCollectDateTimeRequest() *CollectDateTimeRequests {
	return &CollectDateTimeRequests{
		Input:     make(chan *LogLineStruct),
		CountData: make(map[Date]HoursHits),
	}
}

// Run runs the an infinity loop for make things with data
func (c *CollectDateTimeRequests) Run() {

	var line *LogLineStruct
	var ok bool

	for {

		line, ok = <-c.Input
		if !ok {
			return
		}

		date := Date(line.Date.Format("2006-01-02")) // cast to Date format
		if _, ok := c.CountData[date]; !ok {
			c.CountData[date] = make(HoursHits)
		}
		c.CountData[date][line.Date.Format("15")]++

	}

}

// GetChannel returns the channel where send data
func (c *CollectDateTimeRequests) GetChannel() chan *LogLineStruct {
	return c.Input
}