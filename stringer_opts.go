package sloppy_duration

import (
	"text/template"
	"time"
)

type StringerOpts struct {
	MinimumThreshold time.Duration
	MinimumString    string
	Template         *template.Template

	MillisecondThreshold time.Duration
	SecondThreshold      time.Duration
	MinuteThreshold      time.Duration
	HourThreshold        time.Duration
	DayThreshold         time.Duration
	WeekThreshold        time.Duration
	MonthThreshold       time.Duration
}

func (so *StringerOpts) applyDefaults() *StringerOpts {
	so = &(*so)
	def := DefaultStringerOpts()

	if so.MinimumThreshold == 0 {
		so.MinimumThreshold = def.MinimumThreshold
	}
	if so.MinimumString == "" {
		so.MinimumString = def.MinimumString
	}
	if so.Template == nil {
		so.Template = def.Template
	}

	if so.MillisecondThreshold <= 0 {
		so.MillisecondThreshold = def.MillisecondThreshold
	}
	if so.SecondThreshold <= 0 {
		so.SecondThreshold = def.SecondThreshold
	}
	if so.MinuteThreshold <= 0 {
		so.MinuteThreshold = def.MinuteThreshold
	}
	if so.HourThreshold <= 0 {
		so.HourThreshold = def.HourThreshold
	}
	if so.DayThreshold <= 0 {
		so.DayThreshold = def.DayThreshold
	}
	if so.WeekThreshold <= 0 {
		so.WeekThreshold = def.WeekThreshold
	}
	if so.MonthThreshold <= 0 {
		so.MonthThreshold = def.MonthThreshold
	}
	return so
}

func DefaultStringerOpts() *StringerOpts {
	return &StringerOpts{
		MinimumThreshold: time.Second,
		MinimumString:    "< 1s",
		Template:         mustParseTemplate("{{.Value}}{{.ShortUnit}}"),

		MillisecondThreshold: time.Second,
		SecondThreshold:      time.Minute,
		MinuteThreshold:      time.Hour,
		HourThreshold:        Day,
		DayThreshold:         Week * 2,
		WeekThreshold:        Month,
		MonthThreshold:       Year,
	}
}

func mustParseTemplate(str string) *template.Template {
	tpl, err := template.New("sloppy_duration").Parse(str)
	if err != nil {
		panic(err)
	}
	return tpl
}
