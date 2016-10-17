package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type (

	// MapRegexp maps a regexp with params
	MapRegexp map[string]string

	// LogLineStruct keeps some infos of log line
	LogLineStruct struct {
		IP     string
		Date   time.Time
		Method string
		URL    string
		Status int
	}
)

// newMapRegexp map a regexp with it's params
func newMapRegexp(matches []string, param []string) (MapRegexp, error) {

	if len(matches) == 0 {
		return nil, errors.New("row don't match regexp")
	}

	if len(matches) != len(param) {
		return nil, errors.New("mismatch len array")
	}

	splitMatches := make(MapRegexp)
	for i := 1; i < len(param); i++ {
		splitMatches[param[i]] = matches[i]
	}

	return splitMatches, nil

}

// NewLogLineStruct create a new LogLineStruct from a slice of string
func NewLogLineStruct(matches []string, param []string) (*LogLineStruct, error) {

	mapr, err := newMapRegexp(matches, param)
	if err != nil {
		return nil, err
	}

	// find the type of the request
	methodURL := strings.Split(mapr["request"], " ")

	if len(methodURL) != 3 {
		return nil, errors.New("i want 3 spaces")
	}

	// convert the status
	status, err := strconv.Atoi(mapr["status"])
	if err != nil {
		return nil, err
	}

	// parse the time and create a new time.Time struct
	date, err := time.Parse("02/Jan/2006:15:04:05 -0700", mapr["time_local"])
	if err != nil {
		return nil, err
	}

	return &LogLineStruct{
		IP:     mapr["remote_addr"],
		Date:   date,
		Method: methodURL[0],
		URL:    methodURL[1],
		Status: status,
	}, nil

}
