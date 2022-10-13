package currency

type CurrencyCreateRequest struct {
	Title string `json:"title" validate:"required"`
	IsoCode string `json:"iso_code" validate:"required"`
}
