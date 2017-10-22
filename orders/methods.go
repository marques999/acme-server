package orders

func getQueryOptions(orderId string, customerId string) map[string]interface{} {
	if customerId == "admin" {
		return map[string]interface{}{
			"id": orderId,
		}
	} else {
		return map[string]interface{}{
			"id":       orderId,
			"customer": customerId,
		}
	}
}
