package measure

import "fmt"

type Gigabytes float64

// string such as "10 MB" or "5 GB" to gigs
func GigabytesFromString(size string) Gigabytes {
	return MegabytesFromString(size).ToGigabytes()
}

func (gb Gigabytes) ToMegabytes() Megabytes {
	return Megabytes(gb * 1024)
}

func (gb Gigabytes) String() string {
	return fmt.Sprintf("%.2f", gb)
}
