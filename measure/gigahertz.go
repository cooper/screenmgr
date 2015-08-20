package measure

import "fmt"

type Gigahertz float64

func GigahertzFromString(freq string) Gigahertz {
	return MegahertzFromString(freq).ToGigahertz()
}

func (ghz Gigahertz) ToMegahertz() Megahertz {
	return Megahertz(ghz * 1000)
}

func (ghz Gigahertz) String() string {
	return fmt.Sprintf("%.2f", ghz)
}
