package helpers

func GetTableFromType(companyType string) string {
	if companyType == "s" {
		return "shipment_companies"
	}
	if companyType == "c" {
		return "commodity_companies"
	}
	return ""
}
