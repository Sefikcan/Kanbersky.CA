package currency

type CurrencyListResponse struct {
	TotalCount int64 `json:"total_count"`
	TotalPages int `json:"total_pages"`
	Page int `json:"page"`
	Limit int `json:"limit"`
	Currencies []*CurrencyResponse `json:"currencies"`
}
