package measure

import (
	"fmt"
	"strconv"
	"strings"
)

type Megahertz float64

func MegahertzFromString(freq string) Megahertz {
	freq = strings.Replace(strings.TrimSuffix(freq, "Hz"), " ", "", -1)

	// parse the numeric part
	value, err := strconv.ParseFloat(freq[:len(freq)-1], 64)
	if err != nil {
		return 0
	}

	// multipliers
	var multiplier float64 = 1
	switch freq[len(freq)-1:] {
	case "K":
		multiplier *= 1 / 1000
	case "G":
		multiplier *= 1000
	case "T":
		multiplier *= 1000 ^ 2
	}

	return Megahertz(value * multiplier)
}

func (mhz Megahertz) ToGigahertz() Gigahertz {
	return Gigahertz(mhz / 1000)
}

func (mhz Megahertz) String() string {
	return fmt.Sprintf("%.2f", mhz)
}
