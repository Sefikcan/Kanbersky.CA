package currency

type CurrencyUpdateRequest struct {
	ID int `json:"id"`
	Title string `json:"title"`
	IsoCode string `json:"iso_code"`
}
