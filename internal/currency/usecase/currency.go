package usecase

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/sefikcan/kanbersky.ca/internal/currency/mapping"
	"github.com/sefikcan/kanbersky.ca/internal/currency/repository"
	request "github.com/sefikcan/kanbersky.ca/internal/dto/request/currency"
	response "github.com/sefikcan/kanbersky.ca/internal/dto/response/currency"
	"github.com/sefikcan/kanbersky.ca/pkg/config"
	"github.com/sefikcan/kanbersky.ca/pkg/logger"
	"github.com/sefikcan/kanbersky.ca/pkg/util"
	"net/http"
)

type CurrencyUseCase interface {
	Create(ctx context.Context, request *request.CurrencyCreateRequest) (*response.CurrencyResponse, error)
	Update(ctx context.Context, request *request.CurrencyUpdateRequest) (*response.CurrencyResponse, error)
	GetById(ctx context.Context, id int) (*response.CurrencyResponse, error)
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context, request *request.CurrencyPageableRequest) (response.CurrencyListResponse, error)
}

type currencyUseCase struct {
	cfg *config.Config
	currencyRepository repository.CurrencyRepository
	currencyRedisRepository repository.CurrencyRedisRepository
	logger logger.Logger
}

func (c currencyUseCase) Create(ctx context.Context, request *request.CurrencyCreateRequest) (*response.CurrencyResponse, error) {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "currencyUseCase.Create")
	defer span.Finish()

	if err := util.ValidateStruct(&request); err != nil {
		return nil, util.NewHttpResponse(http.StatusBadRequest, util.BadRequest.Error() , errors.WithMessage(err,"currencyUseCase.Create.ValidateStruct"))
	}

	currency := mapping.CreateMapEntity(request)

	resp, err := c.currencyRepository.Create(spanContext, currency)
	if err != nil {
		return nil, err
	}

	mappedResponse := mapping.MapDto(resp)

	if err := c.currencyRedisRepository.Set(spanContext, fmt.Sprintf("%s: %v", "currency", mappedResponse.ID), 3600, mappedResponse); err != nil {
		c.logger.Errorf("currencyUseCase.Create.SetCache: %s", err)
	}

	return mappedResponse, nil
}

func (c currencyUseCase) Update(ctx context.Context, request *request.CurrencyUpdateRequest) (*response.CurrencyResponse, error) {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "currencyUseCase.Update")
	defer span.Finish()

	if err := util.ValidateStruct(&request); err != nil {
		return nil, util.NewHttpResponse(http.StatusBadRequest, util.BadRequest.Error() , errors.WithMessage(err,"currencyUseCase.Create.ValidateStruct"))
	}

	_, err := c.currencyRepository.GetById(spanContext, request.ID)
	if err != nil {
		return nil, err
	}

	currency := mapping.UpdateMapEntity(request)

	updatedCurrency, err := c.currencyRepository.Update(spanContext, currency)
	if err != nil {
		return nil, err
	}

	mappedResponse := mapping.MapDto(updatedCurrency)

	if err := c.currencyRedisRepository.Set(spanContext, fmt.Sprintf("%s: %v", "currency", mappedResponse.ID),  3600, mappedResponse); err != nil {
		c.logger.Errorf("currencyUseCase.Update.SetCache: %s", err)
	}

	return mappedResponse, nil
}

func (c currencyUseCase) GetById(ctx context.Context, id int) (*response.CurrencyResponse, error) {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "currencyUseCase.GetById")
	defer span.Finish()

	currency, err := c.currencyRedisRepository.GetByKey(spanContext, fmt.Sprintf("%s: %v", "currency", id))
	if err != nil {
		c.logger.Errorf("currencyUseCase.GetById.Redis: %v", err)
	} else {
		return currency, nil
	}

	currentCurrency, err := c.currencyRepository.GetById(spanContext, id)
	if err != nil {
		return nil, err
	}

	mappedResponse := mapping.MapDto(currentCurrency)

	if err := c.currencyRedisRepository.Set(spanContext, fmt.Sprintf("%s: %v", "currency", mappedResponse.ID), 3600, mappedResponse); err != nil {
		c.logger.Errorf("currencyUseCase.GetById.SetCache: %s", err)
	}

	return mappedResponse, nil
}

func (c currencyUseCase) Delete(ctx context.Context, id int) error {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "currencyUseCase.Delete")
	defer span.Finish()

	_, err := c.currencyRepository.GetById(spanContext, id)
	if err != nil {
		return err
	}

	if err = c.currencyRepository.Delete(spanContext, id); err != nil {
		return err
	}

	if err = c.currencyRedisRepository.Delete(spanContext, fmt.Sprintf("%s: %v", "currency", id)); err != nil {
		c.logger.Errorf("currencyUseCase.Delete.DeleteCache: %s", err)
	}

	return nil
}

func (c currencyUseCase) GetAll(ctx context.Context, pageableRequest *request.CurrencyPageableRequest) (response.CurrencyListResponse, error) {
	span, spanContext := opentracing.StartSpanFromContext(ctx, "currencyUseCase.GetAll")
	defer span.Finish()

	totalCount := c.currencyRepository.GetCount(spanContext)
	if totalCount == 0 {
		return response.CurrencyListResponse{
			TotalCount: totalCount,
			TotalPages:util.GetTotalPages(totalCount, pageableRequest.Size),
			Page: pageableRequest.Page,
			Limit: pageableRequest.Size,
			Currencies: make([]*response.CurrencyResponse,0),
		}, nil
	}

	var pagination = util.Pagination{
		Page:  pageableRequest.Page,
		Limit: pageableRequest.Size,
	}

	currencies := c.currencyRepository.GetAll(spanContext, pagination)

	return response.CurrencyListResponse{
		TotalCount: totalCount,
		TotalPages:util.GetTotalPages(totalCount, pageableRequest.Size),
		Page: pageableRequest.Page,
		Limit: pageableRequest.Size,
		Currencies: mapping.MapListDto(currencies),
	}, nil
}

func NewCurrencyUseCase(cfg *config.Config, currencyRepository repository.CurrencyRepository, currencyRedisRepository repository.CurrencyRedisRepository, logger logger.Logger) CurrencyUseCase {
	return &currencyUseCase{
		cfg: cfg,
		currencyRepository: currencyRepository,
		currencyRedisRepository: currencyRedisRepository,
		logger: logger,
	}
}