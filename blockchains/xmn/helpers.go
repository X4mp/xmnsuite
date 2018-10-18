package xmn

import "strconv"

func fetchParam(data map[string]string, keyname string) string {
	if val, ok := data[keyname]; ok {
		return val
	}

	return ""
}

func fetchParamWithDefault(data map[string]string, keyname string, defaultValue string) string {
	val := fetchParam(data, keyname)
	if val == "" {
		return defaultValue
	}

	return val
}

func fetchIndex(data map[string]string) int {
	indexAsString := fetchParamWithDefault(data, "index", "0")
	index, indexErr := strconv.Atoi(indexAsString)
	if indexErr != nil {
		index = 0
	}

	if index < 0 {
		index = 0
	}

	return index
}

func fetchAmount(data map[string]string) int {
	amountAsString := fetchParamWithDefault(data, "amount", "20")
	amount, amountErr := strconv.Atoi(amountAsString)
	if amountErr != nil {
		amount = 20
	}

	if amount > 100 {
		amount = 0
	}

	return amount
}
