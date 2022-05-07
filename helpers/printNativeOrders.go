package helpers

import (
	"fmt"
)

func PrintNativeOrders(orders [](*Order)) {
	result := ""

	for _, o := range orders {
		isPref := "⬛️"
		if o.PrefBroker.Name != "" && o.PrefBroker.Surname != "" {
			isPref = "⭐️"
		}

		if o.State == "executed" {
			isPref = "✅"
		}

		result += fmt.Sprintf(
			"💼%s %d. Owner: %s %s (%s)\n\t%s %d %s of %s (%s)\n",
			isPref,
			o.Id,
			o.Owner.Name,
			o.Owner.Surname,
			o.Owner.Email,
			o.Side,
			o.Commodity.Volume,
			o.Commodity.Unit,
			o.Commodity.Label,
			o.State,
		)
	}

	fmt.Print(result)
}
