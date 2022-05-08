package helpers

import (
	"fmt"
)

func PrintForeignOrders(orders [](*Order)) {
	result := ""

	for _, o := range orders {
		result += fmt.Sprintf(
			"💼 %d %s. Owner: %s %s (%s) %s %f %s of %s\n",
			o.Id,
			o.Owner.ExchangerTag,
			o.Owner.Name,
			o.Owner.Surname,
			o.Owner.Email,
			o.Side,
			o.Commodity.Volume,
			o.Commodity.Unit,
			o.Commodity.Label,
		)
	}

	fmt.Print(result)
}
