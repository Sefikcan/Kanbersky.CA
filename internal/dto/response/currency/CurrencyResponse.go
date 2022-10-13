package currency

type CurrencyResponse struct {
	ID int `json:"id"`
	Title string `json:"title"`
	IsoCode string `json:"iso_code"`
}
