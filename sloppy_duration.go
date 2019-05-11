package sloppy_duration

import (
	"bytes"
	"errors"
	"strconv"
	"time"
)

const (
	Day   = time.Hour * 24
	Week  = Day * 7
	Month = Year / 12
	Year  = Day * 365
)

const (
	DayUnit   = "d"
	WeekUnit  = "w"
	MonthUnit = "M"
	YearUnit  = "y"
)


type stringerData struct {
	Value     int
	ShortUnit string
	LongUnit  string
}

// A sloppy wrapper around a time.Duration
type SloppyDuration struct {
	dur   time.Duration
	sOpts *StringerOpts
}

// Wrap the given time.Duration so it can be manipulated sloppily.
// Uses default StringerOpts.
func Wrap(dur time.Duration) *SloppyDuration {
	return WrapWithOpts(dur, DefaultStringerOpts())
}

// Wrap the given time.Duration so it can be manipulated sloppily.
// Use the supplied StringerOpts to customize the output of String().
func WrapWithOpts(dur time.Duration, opts *StringerOpts) *SloppyDuration {
	return &SloppyDuration{
		dur:   dur,
		sOpts: opts.applyDefaults(),
	}
}

// Parse a sloppy time duration string. The units "d" (day), "w" (week),
// "M" (month), and "y" (year) are supported in addition to those supported
// by time.Duration. Multi-unit durations, such as "3h1m30s" are not supported
// (because how is that sloppy?). Signed durations are not supported.
// Uses default StringerOpts.
func Parse(str string) (*SloppyDuration, error) {
	return ParseWithOpts(str, DefaultStringerOpts())
}

// Same as Parse() but use the supplied StringerOpts to customize the output of String().
func ParseWithOpts(str string, opts *StringerOpts) (*SloppyDuration, error) {
	if str == "" {
		return nil, errors.New("empty string")
	}

	// Told you it was sloppy

	unit := lastNonDigits(str)
	valStr := str[:len(str)-len(unit)]

	// Ensure non-multi-part (i.e 1m30s)
	_, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return nil, err
	}

	dAsNano, err := time.ParseDuration(valStr + "ns")
	if err != nil {
		return nil, err
	}

	var d time.Duration
	var denomination time.Duration

	switch unit {
	case DayUnit:
		denomination = Day
	case WeekUnit:
		denomination = Week
	case MonthUnit:
		denomination = Month
	case YearUnit:
		denomination = Year
	default:
		denomination = 0
		d, err = time.ParseDuration(str)
		if err != nil {
			return nil, err
		}
	}

	if denomination != 0 {
		d = dAsNano * denomination
	}

	return &SloppyDuration{
		dur:   d,
		sOpts: opts.applyDefaults(),
	}, nil
}

// Get the wrapped time.Duration value
func (sd *SloppyDuration) Duration() time.Duration {
	return sd.dur
}

// Get the rough number of milliseconds in the duration
func (sd *SloppyDuration) Milliseconds() int {
	return int(sd.dur.Nanoseconds() / 1e+6)
}

// Get the rough number of seconds in the duration
func (sd *SloppyDuration) Seconds() int {
	return int(sd.dur.Seconds())
}

// Get the rough number of minutes in the duration
func (sd *SloppyDuration) Minutes() int {
	return int(sd.dur.Minutes())
}

// Get the rough number of hours in the duration
func (sd *SloppyDuration) Hours() int {
	return int(sd.dur.Hours())
}

// Get the rough number of days in the duration
func (sd *SloppyDuration) Days() int {
	days := sd.dur / Day
	return int(days)
}

// Get the rough number of weeks in the duration
func (sd *SloppyDuration) Weeks() int {
	weeks := sd.dur / Week
	return int(weeks)
}

// Get the rough number of months in the duration
func (sd *SloppyDuration) Months() int {
	months := sd.dur / Month
	return int(months)
}

// Get the rough number of years in the duration
func (sd *SloppyDuration) Years() int {
	years := sd.dur / Year
	return int(years)
}

// Generate a sloppy string representation of the duration
func (sd *SloppyDuration) String() string {

	var sData *stringerData

	switch {
	case sd.dur < sd.sOpts.MinimumThreshold:
		return sd.sOpts.MinimumString

	case sd.dur < sd.sOpts.MillisecondThreshold:
		val := sd.Milliseconds()
		sData = &stringerData{
			Value:     val,
			ShortUnit: "ms",
			LongUnit:  pluralize("millisecond", val),
		}
	case sd.dur < sd.sOpts.SecondThreshold:
		val := sd.Seconds()
		sData = &stringerData{
			Value:     val,
			ShortUnit: "s",
			LongUnit:  pluralize("second", val),
		}
	case sd.dur < sd.sOpts.MinuteThreshold:
		val := sd.Minutes()
		sData = &stringerData{
			Value:     val,
			ShortUnit: "m",
			LongUnit:  pluralize("minute", val),
		}
	case sd.dur < sd.sOpts.HourThreshold:
		val := sd.Hours()
		sData = &stringerData{
			Value:     val,
			ShortUnit: "h",
			LongUnit:  pluralize("hour", val),
		}
	case sd.dur < sd.sOpts.DayThreshold:
		val := sd.Days()
		sData = &stringerData{
			Value:     val,
			ShortUnit: DayUnit,
			LongUnit:  pluralize("day", val),
		}
	case sd.dur < sd.sOpts.WeekThreshold:
		val := sd.Weeks()
		sData = &stringerData{
			Value:     val,
			ShortUnit: WeekUnit,
			LongUnit:  pluralize("week", val),
		}
	case sd.dur < sd.sOpts.MonthThreshold:
		val := sd.Months()
		sData = &stringerData{
			Value:     val,
			ShortUnit: MonthUnit,
			LongUnit:  pluralize("month", val),
		}
	default:
		val := sd.Years()
		sData = &stringerData{
			Value:     val,
			ShortUnit: YearUnit,
			LongUnit:  pluralize("year", val),
		}
	}

	out := bytes.NewBuffer([]byte{})
	err := sd.sOpts.Template.Execute(out, sData)
	if err != nil {
		panic(err)
	}

	return out.String()
}

func pluralize(str string, value int) string {
	if value > 1 {
		return str + "s"
	}
	return str
}

func lastNonDigits(str string) string {
	result := make([]byte, 0)

	for i := len(str)-1; i>=0; i-- {
		ch := str[i]
		if ch < '0' || ch > '9' {
			result = append([]byte{ch}, result...)
		}
	}

	return string(result)
}
