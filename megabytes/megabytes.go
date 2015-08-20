package megabytes

import (
	"github.com/pivotal-golang/bytefmt"
	"strings"
)

// string such as "10 MB" or "1 GB" to numeric mb
func MegabytesFromString(size string) float64 {
	size = strings.Replace(strings.TrimSuffix(size, "B"), " ", "", -1)
	mb, err := bytefmt.ToMegabytes(size)
	if err != nil {
		return 0
	}
	return float64(mb)
}

// this seems pretty useless, but it may be necessary in case
// it becomes useful to use real IS units
func MegabytesToGigabytes(mb float64) float64 {
	return mb / 1024
}

// string such as "10 MB" or "5 GB" to numeric gigs
func GigabytesFromString(size string) float64 {
	return MegabytesToGigabytes(MegabytesFromString(size))
}
