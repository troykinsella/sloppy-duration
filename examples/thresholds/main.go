package main

import (
	"fmt"
	"github.com/troykinsella/sloppy-duration"
	"time"
)

func main() {

	// Show seconds until > 90s, above which show minutes
	ninetySeconds, err := time.ParseDuration("90s")
	if err != nil {
		panic(err)
	}
	sOpts := &sloppy_duration.StringerOpts{
		SecondThreshold: ninetySeconds,
	}

	d, err := sloppy_duration.ParseWithOpts("80s", sOpts)
	if err != nil {
		panic(err)
	}
	fmt.Println(d)
	// Prints "80s"

	d, err = sloppy_duration.ParseWithOpts("100s", sOpts)
	if err != nil {
		panic(err)
	}
	fmt.Println(d)
	// Prints "1m"
}

