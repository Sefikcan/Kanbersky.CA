package repository

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sefikcan/kanbersky.ca/internal/currency/entity"
	"github.com/sefikcan/kanbersky.ca/pkg/util"
	"gorm.io/gorm"
)

type CurrencyRepository interface {
	Create(ctx context.Context, currency entity.Currency) (entity.Currency, error)
	Update(ctx context.Context, currency entity.Currency) (entity.Currency, error)
	GetById(ctx context.Context, id int) (entity.Currency, error)
	Delete(ctx context.Context, id int) error
	GetCount(ctx context.Context) int64
	GetAll(ctx context.Context, query util.Pagination) []entity.Currency
}

type currencyRepository struct {
	db *gorm.DB
}

func (c currencyRepository) Create(ctx context.Context, currency entity.Currency) (entity.Currency, error) {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "currencyRepository.Create")
	defer span.Finish()

	if result := c.db.WithContext(spanContext).Create(&currency); result.Error != nil {
		return entity.Currency{}, errors.Wrap(result.Error, "currencyRepository.Create.DbError")
	}

	return currency, nil
}

func (c currencyRepository) Update(ctx context.Context, currency entity.Currency) (entity.Currency, error) {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "currencyRepository.Update")
	defer span.Finish()

	if result := c.db.WithContext(spanContext).Save(&currency); result.Error != nil {
		return entity.Currency{}, errors.Wrap(result.Error, "currencyRepository.Update.DbError")
	}

	return currency, nil
}

func (c currencyRepository) GetById(ctx context.Context, id int) (entity.Currency, error) {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "currencyRepository.GetById")
	defer span.Finish()

	currentCurrency := entity.Currency{}
	err := c.db.WithContext(spanContext).Where(`id = ?`, id).First(&currentCurrency).Error
	if err != nil {
		return entity.Currency{}, errors.Wrap(err,"currencyRepository.GetById.DbError")
	}

	return currentCurrency, err
}

func (c currencyRepository) Delete(ctx context.Context, id int) error {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "currencyRepository.Delete")
	defer span.Finish()

	if result := c.db.WithContext(spanContext).Delete(&entity.Currency{ID: id}); result.Error != nil {
		return errors.Wrap(result.Error, "currencyRepository.Delete.DbError")
	}

	return nil
}

func (c currencyRepository) GetCount(ctx context.Context) int64 {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "currencyRepository.GetCount")
	defer span.Finish()

	var currencies []*entity.Currency

	var totalCount int64
	c.db.WithContext(spanContext).Model(currencies).Count(&totalCount)

	return totalCount
}

func (c currencyRepository) GetAll(ctx context.Context, query util.Pagination) []entity.Currency {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "currencyRepository.GetAll")
	defer span.Finish()

	var currencies []entity.Currency
	c.db.WithContext(spanContext).Offset(query.GetOffset()).Limit(query.GetLimit()).Order(query.GetLimit()).Find(&currencies)

	return currencies
}

func NewCurrencyRepository(db *gorm.DB) CurrencyRepository {
	return &currencyRepository{
		db: db,
	}
}