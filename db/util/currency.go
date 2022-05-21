package util

const (
	// Supported Currency
	USD = "USD"
	EUR = "EUR"
	GBP = "GBP"
	HKD = "HKD"
	CAD = "CAD"
	JPY = "JPY"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, GBP, HKD, CAD, JPY:
		return true
	}
	return false
}
