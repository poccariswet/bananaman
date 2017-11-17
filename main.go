package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	var ID, start, radiname string
	flag.StringVar(&ID, "id", "", "id")
	flag.StringVar(&start, "s", "", "start")
	flag.StringVar(&radiname, "file", "", "file")
	flag.Parse()
	if ID == "" {
		fmt.Println("Please input stationID, like '-id=TBS'")
		os.Exit(1)
	}
	if start == "" {
		fmt.Printf("Please input start time you wanna listen to radio name,\nlike 2017/11/11/01:00 -> 20171111010000\n")
		os.Exit(1)
	}
	if radiname == "" {
		fmt.Println("Please input filename,\nlike '-file=bananamoonGOLD'")
	}
	go spinner(100 * time.Millisecond)
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}
