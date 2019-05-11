# Golang `SloppyDuration`

A Golang library for manipulating durations with coarser granularity and less precision.

## Why?

The `time.Duration` type supports a time granularity of up to an "hour" unit,
and it does so to ensure duration precision. Obviously, for example,
a month is a variable duration, which is why it cannot be manipulated precisely.

`SloppyDuration` expands upon `time.Duration`'s parsing and `Stringer` ability 
with support for:
* Days (`d`): `24h`
* Weeks (`w`): `24h * 7`
* Months (`M`): `365d / 12`
* Years (`y`): `1d * 365`

Additionally, `String()` provides a simpler, less precise output, such as
"2m" rather than "2m3s". The string value can also be customized via a `template.Template`.

`SloppyDuration` does not support multi-unit duration strings, such as "1h45m0s",
nor signed durations.

An example use case for this library could be to render a publication
date for a blog post of, say, "2 months ago".

## Usage

Install:

```bash
go get github.com/troykinsella/sloppy-duration
```

Use:

```go
year, err := sloppy_duration.Parse("1y")
if err != nil {
	panic(err)
}
fmt.Printf("%s days in a year\n", year.Days())
// Prints "365 days in a year"

soon, err := sloppy_duration.Parse("48h")
if err != nil {
	panic(err)
}
fmt.Printf("Production will fail in: %s\n", soon)
// Prints "Production will fail in: 2d"
```

See `examples/` for more details.

## Testing

```bash
go get github.com/onsi/ginkgo/ginkgo
ginkgo
```

## License

MIT © Troy Kinsella
