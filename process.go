package main

import (
	"log"
	"regexp"
	"sync"
)

// FromLineToStruct builds a LogLineStruct or print an error
func FromLineToStruct(input chan *RowLine, reLine regexp.Regexp, wait *sync.WaitGroup) {

	defer wait.Done()

	var line *RowLine
	var ok bool

	for {

		line, ok = <-input

		// if !ok, channel has been closed, so we can return
		if !ok {
			return
		}

		structLine, err := NewLogLineStruct(reLine.FindStringSubmatch(line.RowStr), reLine.SubexpNames())
		if err != nil {
			log.Println(line.Num, " - ", err)
		} else {
			for _, chanel := range GlobalCollectors {
				// we are going to pass a new struct!!
				chanel.GetChannel() <- *structLine
			}
		}

	}

}
