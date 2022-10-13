package mapping

import (
	"github.com/sefikcan/kanbersky.ca/internal/currency/entity"
	"github.com/sefikcan/kanbersky.ca/internal/dto/request/currency"
)

func CreateMapEntity(currency *currency.CurrencyCreateRequest) entity.Currency {
	return entity.Currency{
		Title: currency.Title,
		IsoCode: currency.IsoCode,
	}
}

func UpdateMapEntity(currency *currency.CurrencyUpdateRequest) entity.Currency {
	return entity.Currency{
		ID: currency.ID,
		Title: currency.Title,
		IsoCode: currency.IsoCode,
	}
}
