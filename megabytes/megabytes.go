package megabytes

import (
	"fmt"
	"github.com/pivotal-golang/bytefmt"
	"strings"
)

type Megabytes float64

// string such as "10 MB" or "1 GB" to numeric mb
func MegabytesFromString(size string) Megabytes {
	size = strings.Replace(strings.TrimSuffix(size, "B"), " ", "", -1)
	mb, err := bytefmt.ToMegabytes(size)
	if err != nil {
		return 0
	}
	return Megabytes(mb)
}

func (mb Megabytes) ToGigabytes() Gigabytes {
	return Gigabytes(mb / 1024)
}

func (mb Megabytes) String() string {
	return fmt.Sprintf("%.2f", mb)
}
