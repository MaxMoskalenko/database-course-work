package helpers

import (
	"fmt"
)

func PrintPersonalOrders(orders [](*Order)) {
	result := ""

	for _, o := range orders {
		prefBroker := ""
		if o.PrefBroker.Name != "" && o.PrefBroker.Surname != "" {
			prefBroker = fmt.Sprintf("Selected broker: %s %s", o.PrefBroker.Name, o.PrefBroker.Surname)
		}

		result += fmt.Sprintf(
			"💼 %d. %s %f %s of %s (%s) %s\n",
			o.Id,
			o.Side,
			o.Commodity.Volume,
			o.Commodity.Unit,
			o.Commodity.Label,
			o.State,
			prefBroker,
		)
	}

	fmt.Print(result)
}
