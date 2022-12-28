package services

var data map[string]int = map[string]int{}

func SetIgnoredCurrency(symbol string, value int) {
	data[symbol] = value
}

func GetIgnoredCurrencies() *[]string {
	var symbols []string = []string{}

	for key, val := range data {
		symbols = append(symbols, key)

		if val > 1 {
			data[key] = val - 1
		} else {
			delete(data, key)
		}
	}

	if len(symbols) == 0 {
		return nil
	}

	return &symbols
}
