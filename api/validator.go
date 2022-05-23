package api

import (
	"github.com/go-playground/validator/v10"
	"simplebank/db/util"
)

var validCurrency validator.Func = func(fleidLevel validator.FieldLevel) bool {
	if currency, ok := fleidLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}
