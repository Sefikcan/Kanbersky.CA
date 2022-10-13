package mapping

import (
	"github.com/sefikcan/kanbersky.ca/internal/currency/entity"
	"github.com/sefikcan/kanbersky.ca/internal/dto/response/currency"
)

type CurrencyResponses []*currency.CurrencyResponse

func MapDto(c entity.Currency) *currency.CurrencyResponse {
	return &currency.CurrencyResponse{
		ID: c.ID,
		Title: c.Title,
		IsoCode: c.IsoCode,
	}
}

func MapListDto(currencies []entity.Currency) CurrencyResponses {
	var currencyResp CurrencyResponses
	for _, c := range currencies {
		mappedCurrency := MapDto(c)
		currencyResp = append(currencyResp, mappedCurrency)
	}

	return currencyResp
}
