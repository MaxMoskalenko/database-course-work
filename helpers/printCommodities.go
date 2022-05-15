package helpers

import (
	"fmt"
	"sort"
)

func PrintCommodities(commodities [](*Commodity)) {
	sort.SliceStable(commodities, func(i, j int) bool {
		return commodities[i].Owner.Email < commodities[j].Owner.Email
	})
	previousUser := ""
	result := ""

	for _, c := range commodities {
		if c.Owner.Email != previousUser {
			result += fmt.Sprintf("🥸  %s %s`s\n", c.Owner.Name, c.Owner.Surname)
			previousUser = c.Owner.Email
		}
		result += fmt.Sprintf("📦 %f %s of %s\n", c.Volume, c.Unit, c.Label)
	}
	fmt.Print(result)
}
