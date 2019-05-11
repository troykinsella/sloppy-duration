package main

import (
	"fmt"
	"github.com/troykinsella/sloppy-duration"
	"time"
)

func main() {
	oneMinute, err := time.ParseDuration("1m")
	if err != nil {
		panic(err)
	}

	sOpts := &sloppy_duration.StringerOpts{
		MinimumThreshold: oneMinute,
		MinimumString: "less than a minute",
	}

	d, err := sloppy_duration.ParseWithOpts("59s", sOpts)
	if err != nil {
		panic(err)
	}
	fmt.Println(d)
	// Prints "less than a minute"

	d, err = sloppy_duration.ParseWithOpts("61s", sOpts)
	if err != nil {
		panic(err)
	}
	fmt.Println(d)
	// Prints "1m"
}
