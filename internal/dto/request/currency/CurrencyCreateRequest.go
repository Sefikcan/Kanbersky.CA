package currency

type CurrencyCreateRequest struct {
	Title string `json:"title" validate:"required,min=3,max=12"`
	IsoCode string `json:"iso_code" validate:"required,min=3"`
}
