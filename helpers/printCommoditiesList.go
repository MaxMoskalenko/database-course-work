package helpers

import (
	"fmt"
)

func PrintCommoditiesList(commodities [](*Commodity)) {
	result := ""

	for _, c := range commodities {
		result += fmt.Sprintf("📦 %s in %ss\n", c.Label, c.Unit)
	}
	fmt.Print(result)
}
