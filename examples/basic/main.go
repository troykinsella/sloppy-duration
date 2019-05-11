package main

import (
	"fmt"
	"github.com/troykinsella/sloppy-duration"
	"time"
)

func main() {

	// Parse a value
	twoDays, err := sloppy_duration.Parse("2d")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Two days is: %s\n", twoDays)
	// Prints "Two days is: 2d"

	// Compare with a time.Duration
	day, err := time.ParseDuration("24h")
	if err != nil {
		panic(err)
	}

	fmt.Printf("2d is 48h: %t\n", twoDays.Duration() == day * 2)
	// Prints "2d is 48h: true"

	// Make a time.Duration sloppy
	normalWeek := day * 7
	fmt.Printf("Normal week: %s\n", normalWeek)
	// Prints: "Normal week: 168h0m0s"
	fmt.Printf("Sloppy week: %s\n", sloppy_duration.Wrap(normalWeek))
	// Prints "Sloppy week: 7d" (The default week unit threshold is >=14 days)
	fmt.Printf("Sloppy couple of weeks: %s\n", sloppy_duration.Wrap(normalWeek * 2))
	// Prints "Sloppy couple of weeks: 2w"

	// Convert units (sloppily)
	year, err := sloppy_duration.Parse("1y")
	if err != nil {
		panic(err)
	}
	fmt.Printf("There are about %d days in a year\n", year.Days())
	// Prints "There are about 365 days in a year"
}
