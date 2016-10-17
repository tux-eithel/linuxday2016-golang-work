package main

import (
	"bufio"
	"errors"
	"flag"
	"log"
	"os"
	"regexp"
	"runtime"
	"sync"
)

const (
	logFormat string = `^(?P<remote_addr>.+) (?P<gzip_ratio>.+) ` +
		`(?P<remote_user>.+) \[(?P<time_local>.+)\] "(?P<request>.+)" ` +
		`(?P<status>.+) (?P<bytes_sent>.+) "(?P<http_referer>.*)" "(?P<http_user_agent>.*)"$`
)

var (
	inputFile string
)

func init() {
	// init the variables
	flag.StringVar(&inputFile, "f", "", "input file to parse")
}

func main() {

	// parse input flags
	flag.Parse()

	// compile the regexp
	reLine := regexp.MustCompile(logFormat)

	// try to open the file
	logFile, err := openFile(inputFile)
	if err != nil {
		log.Fatalln(err)
	}
	// at the end close the file
	defer logFile.Close()

	// a sync.WaitGroup define a counter for a number of goroutines that need to be waited
	waitRoutine := &sync.WaitGroup{}
	waitRoutine.Add(runtime.NumCPU())

	// create a channel to send RowLine
	chRowLine := make(chan *RowLine)

	// how many routine ? for now cpunumber is ok
	for j := 0; j < runtime.NumCPU(); j++ {
		go FromLineToStruct(chRowLine, reLine, waitRoutine)
	}

	// start to scan the file
	logLines := bufio.NewScanner(logFile)
	i := 1
	for logLines.Scan() {

		// send the struct to the channel
		chRowLine <- &RowLine{
			Num:    i,
			RowStr: logLines.Text(),
		}

		i++
		if i == 100 {
			break
		}
	}

	// wait all the goroutine to end
	waitRoutine.Wait()

}

// openFile tries to open a file
func openFile(fileName string) (*os.File, error) {
	if fileName == "" {
		return nil, errors.New("empty file")
	}

	return os.Open(fileName)
}
