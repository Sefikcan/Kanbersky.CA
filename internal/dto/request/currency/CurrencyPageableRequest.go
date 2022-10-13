package currency

type CurrencyPageableRequest struct {
	Size int `json:"size,omitempty"`
	Page int `json:"page,omitempty"`
	OrderBy string `json:"orderBy,omitempty"`
}
