package requests

type CurrencyConfigRequest struct {
	Symbol string `validate:"required" json:"symbol" form:"symbol"`
}
