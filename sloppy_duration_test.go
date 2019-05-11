package sloppy_duration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"text/template"
	"time"
)

var _ = Describe("sloppy_duration", func() {

	Describe("Parse", func() {

		It("parses", func() {
			tests := []struct {
				d        string
				expected string
			}{
				{"1m", "1m"},
				{"2m", "2m"},
				{"1h", "1h"},
				{"2h", "2h"},
				{"1d", "24h"},
				{"2d", "48h"},
				{"7d", "168h"},
				{"1w", "168h"},
				{"4w", "672h"}, // because day * 7 * 4
				{"1M", "730h"}, // because year / 12
				{"2M", "1460h"},
				{"12M", "8760h"},
				{"1y", "8760h"},
				{"2y", "17520h"},
			}

			for _, test := range tests {
				expected, err := time.ParseDuration(test.expected)
				Expect(err).ToNot(HaveOccurred())

				d, err := Parse(test.d)
				Expect(err).ToNot(HaveOccurred())
				Expect(d.dur).To(Equal(expected), "Duration: "+test.d)
			}
		})

		It("errors", func() {
			tests := []struct {
				d   string
				err string
			}{
				{"1c", "time: unknown unit c in duration 1c"},
				{"1m30s", `strconv.ParseFloat: parsing "1m30": invalid syntax`},
				{"turtleh", `strconv.ParseFloat: parsing "turtle": invalid syntax`},
				{"2yy", `strconv.ParseFloat: parsing "2y": invalid syntax`},
			}

			for _, test := range tests {
				_, err := Parse(test.d)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal(test.err), "Duration: "+test.d)
			}
		})
	})

	Describe("Stringer", func() {

		It("customizes unit thresholds", func() {

			tests := []struct{
				d string
				expected string
				opts *StringerOpts
			}{
				{
					"60s",
					"60s",
					&StringerOpts{
						SecondThreshold: 90 * time.Second,
					},
				},
				{
					"89s",
					"89s",
					&StringerOpts{
						SecondThreshold: 90 * time.Second,
					},
				},
				{
					"91s",
					"1m",
					&StringerOpts{
						SecondThreshold: 90 * time.Second,
					},
				},
				{
					"60m",
					"60m",
					&StringerOpts{
						MinuteThreshold: 90 * time.Minute,
					},
				},
				{
					"89m",
					"89m",
					&StringerOpts{
						MinuteThreshold: 90 * time.Minute,
					},
				},
				{
					"91m",
					"1h",
					&StringerOpts{
						MinuteThreshold: 90 * time.Minute,
					},
				},
				{
					"7d",
					"1w", // Would be 7d by default
					&StringerOpts{
						DayThreshold: Week, // Default is 2 weeks
					},
				},
				{
					"364d",
					"364d",
					&StringerOpts{
						DayThreshold: Year,
					},
				},
			}

			for _, test := range tests {
				d, err := ParseWithOpts(test.d, test.opts)
				Expect(err).ToNot(HaveOccurred())
				Expect(d.String()).To(Equal(test.expected), "Duration: "+test.d)
			}

		})

		It("customizes template", func() {
			tests := []struct{
				d string
				expected string
				template string
			}{
				{
					"1d",
					"1 day",
					"{{.Value}} {{.LongUnit}}",
				},
				{
					"2d",
					"2 days",
					"{{.Value}} {{.LongUnit}}",
				},
				{
					"1M",
					"1 month or something",
					"{{.Value}} {{.LongUnit}} or something",
				},
			}

			for _, test := range tests {
				template, err := template.New("test").Parse(test.template)
				Expect(err).ToNot(HaveOccurred())

				d, err := ParseWithOpts(test.d, &StringerOpts{
					Template: template,
				})

				Expect(d.String()).To(Equal(test.expected), "Duration: "+test.d)
			}
		})

		It("customizes minimum threshold", func() {
			tests := []struct{
				d string
				mt string
				ms string
				expected string
			}{
				{
					"500ms",
					"1s",
					"< 1s",
					"< 1s",
				},
				{
					"999ms",
					"1s",
					"< 1s",
					"< 1s",
				},
				{
					"1000ms",
					"1s",
					"< 1s",
					"1s",
				},
				{
					"4h",
					"24h",
					"less than a day",
					"less than a day",
				},
				{
					"25h",
					"24h",
					"under a friggin' day",
					"1d",
				},
			}

			for _, test := range tests {
				mt, err := time.ParseDuration(test.mt)
				Expect(err).ToNot(HaveOccurred())

				d, err := ParseWithOpts(test.d, &StringerOpts{
					MinimumThreshold: mt,
					MinimumString: test.ms,
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(d.String()).To(Equal(test.expected), "Duration: "+test.d)
			}
		})

	})

	Describe("String", func() {

		It("with default stringer options", func() {

			tests := [][]string{
				{"500ms", "< 1s"},
				{"1s", "1s"},
				{"59s", "59s"},
				{"90s", "1m"},
				{"59m", "59m"},
				{"60m", "1h"},
				{"23h", "23h"},
				{"24h", "1d"},
				{"48h", "2d"},
				{"168h", "7d"},
				{"335h", "13d"},
				{"336h", "2w"},
				{"671h", "3w"},
				{"729h", "4w"}, // because day * 7 * 4
				{"730h", "1M"}, // because year / 12
				{"4380h", "6M"},
				{"8759h", "11M"},
				{"8760h", "1y"},
				{"43800h", "5y"},
			}

			for _, test := range tests {
				d, err := time.ParseDuration(test[0])
				Expect(err).ToNot(HaveOccurred())

				sd := Wrap(d)
				str := sd.String()

				Expect(str).To(Equal(test[1]), "Duration: "+test[0])
			}
		})

	})
})
