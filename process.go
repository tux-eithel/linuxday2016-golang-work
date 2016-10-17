package main

import (
	"log"
	"regexp"
	"sync"
)

// FromLineToStruct builds a LogLineStruct or print an error
func FromLineToStruct(input chan *RowLine, reLine *regexp.Regexp, wait *sync.WaitGroup) {

	defer wait.Done()

	var line *RowLine

	for {

		line = <-input

		_, err := NewLogLineStruct(reLine.FindStringSubmatch(line.RowStr), reLine.SubexpNames())
		if err != nil {
			log.Println(line.Num, " - ", err)
		} else {
			// fmt.Println(structLine.Status)
		}

	}

}
