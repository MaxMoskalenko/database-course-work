package helpers

import "fmt"

func PrintRaces(races [](*Race)) {
	for _, r := range races {
		fmt.Printf("%2d. Race of %s: %s-%s on %v\n", r.Id, r.Company.Tag, r.FromExch.Tag, r.ToExch.Tag, r.DateStamp.Format("2006-01-02 15:04"))
	}
}
