package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/sefikcan/kanbersky.ca/internal/currency/usecase"
	"github.com/sefikcan/kanbersky.ca/internal/dto/request/currency"
	"github.com/sefikcan/kanbersky.ca/pkg/config"
	"github.com/sefikcan/kanbersky.ca/pkg/logger"
	"github.com/sefikcan/kanbersky.ca/pkg/util"
	"net/http"
	"strconv"
	"strings"
)

type CurrencyHandlers interface {
	Create() echo.HandlerFunc
	Update() echo.HandlerFunc
	GetById() echo.HandlerFunc
	Delete() echo.HandlerFunc
	GetAll() echo.HandlerFunc
}

type currencyHandlers struct {
	cfg *config.Config
	currencyUseCase usecase.CurrencyUseCase
	logger logger.Logger
}

// Create godoc
// @Summary Create currency
// @Description Create currency handler
// @Tags Currency
// @Accept json
// @Produce json
// @Success 201 {object} currency.CurrencyResponse
// @Router /currencies [post]
func (c currencyHandlers) Create() echo.HandlerFunc {
	return func(e echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(util.GetRequestCtx(e), "currencyHandler.Create")
		defer span.Finish()

		currencyRequest := &currency.CurrencyCreateRequest{}
		if err := e.Bind(currencyRequest); err != nil {
			util.PrepareLogging(e, c.logger, err)
			return e.JSON(http.StatusBadRequest, util.NewHttpResponse(http.StatusBadRequest, strings.ToLower(err.Error()),nil))
		}

		createdCurrency, err := c.currencyUseCase.Create(ctx, currencyRequest)
		if err != nil {
			util.PrepareLogging(e, c.logger, err)
			return e.JSON(http.StatusInternalServerError, util.NewHttpResponse(http.StatusInternalServerError, strings.ToLower(err.Error()),nil))
		}

		return e.JSON(http.StatusCreated, createdCurrency)
	}
}

// Update godoc
// @Summary Update currencies
// @Description Update currency handler
// @Tags Currency
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} currency.CurrencyResponse
// @Router /currencies/{id} [put]
func (c currencyHandlers) Update() echo.HandlerFunc {
	return func(e echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(util.GetRequestCtx(e), "currencyHandler.Update")
		defer span.Finish()

		id, err := strconv.Atoi(e.Param("id"))
		if err != nil {
			util.PrepareLogging(e, c.logger, err)
			return e.JSON(http.StatusBadRequest, util.NewHttpResponse(http.StatusBadRequest, strings.ToLower(err.Error()),nil))
		}

		currency := &currency.CurrencyUpdateRequest{}
		if err = e.Bind(currency); err != nil {
			util.PrepareLogging(e, c.logger, err)
			return e.JSON(http.StatusBadRequest, util.NewHttpResponse(http.StatusBadRequest, strings.ToLower(err.Error()),nil))
		}

		currency.ID = id
		updatedCurrency, err := c.currencyUseCase.Update(ctx, currency)
		if err != nil {
			util.PrepareLogging(e, c.logger, err)
			return e.JSON(http.StatusInternalServerError, util.NewHttpResponse(http.StatusInternalServerError, strings.ToLower(err.Error()),nil))
		}

		return e.JSON(http.StatusOK, updatedCurrency)
	}
}

// GetById godoc
// @Summary Get by id currency
// @Description Get by id currency handler
// @Tags Currencies
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} currency.CurrencyResponse
// @Router /currencies/{id} [get]
func (c currencyHandlers) GetById() echo.HandlerFunc {
	return func(e echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(util.GetRequestCtx(e), "currencyHandler.GetById")
		defer span.Finish()

		id, err := strconv.Atoi(e.Param("id"))
		if err != nil {
			util.PrepareLogging(e, c.logger, err)
			return e.JSON(http.StatusBadRequest, util.NewHttpResponse(http.StatusBadRequest, strings.ToLower(err.Error()),nil))
		}

		currencyCurrency, err := c.currencyUseCase.GetById(ctx, id)
		if err != nil {
			util.PrepareLogging(e, c.logger, err)
			return e.JSON(http.StatusInternalServerError, util.NewHttpResponse(http.StatusInternalServerError, strings.ToLower(err.Error()),nil))
		}

		return e.JSON(http.StatusOK, currencyCurrency)
	}
}

// Delete godoc
// @Summary Delete currency
// @Description Delete by id currency handler
// @Tags Currency
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 204
// @Router /currencies/{id} [delete]
func (c currencyHandlers) Delete() echo.HandlerFunc {
	return func(e echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(util.GetRequestCtx(e), "currencyHandler.Delete")
		defer span.Finish()

		id, err := strconv.Atoi(e.Param("id"))
		if err != nil {
			util.PrepareLogging(e, c.logger, err)
			return e.JSON(http.StatusBadRequest, util.NewHttpResponse(http.StatusBadRequest, strings.ToLower(err.Error()),nil))
		}

		if err = c.currencyUseCase.Delete(ctx, id); err != nil {
			util.PrepareLogging(e, c.logger, err)
			return e.JSON(http.StatusInternalServerError, util.NewHttpResponse(http.StatusInternalServerError, strings.ToLower(err.Error()),nil))
		}

		return e.NoContent(http.StatusNoContent)
	}
}

// GetAll godoc
// @Summary Get all currencies
// @Description Get all currencies with pagination
// @Tags Currencies
// @Accept json
// @Produce json
// @Param page query int false "page number" Format(page)
// @Param size query int false "number of elements per page" Format(size)
// @Param orderBy query int false "filter name" Format(orderBy)
// @Success 200 {object} currency.CurrencyListResponse
// @Router /currencies/ [get]
func (c currencyHandlers) GetAll() echo.HandlerFunc {
	return func(e echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(util.GetRequestCtx(e), "currencyHandler.GetAll")
		defer span.Finish()

		var currencyPageableRequest currency.CurrencyPageableRequest
		if e.QueryParam("page") != "" {
			resp, err := strconv.Atoi(e.QueryParam("page"))
			if err == nil {
				currencyPageableRequest.Page = resp
			}
		} else {
			currencyPageableRequest.Page = 1
		}

		if e.QueryParam("limit") != "" {
			resp, err := strconv.Atoi(e.QueryParam("limit"))
			if err == nil {
				currencyPageableRequest.Size = resp
			}
		} else {
			currencyPageableRequest.Size = 10
		}

		if e.QueryParam("sort") != "" {
			currencyPageableRequest.OrderBy = e.QueryParam("sort")
		}

		currencyList, err := c.currencyUseCase.GetAll(ctx, &currencyPageableRequest)
		if err != nil {
			util.PrepareLogging(e, c.logger, err)
			return e.JSON(http.StatusInternalServerError, util.NewHttpResponse(http.StatusInternalServerError, strings.ToLower(err.Error()),nil))
		}

		return e.JSON(http.StatusOK, currencyList)
	}
}

func NewCurrencyHandler(cfg *config.Config, currencyUseCase usecase.CurrencyUseCase, logger logger.Logger) CurrencyHandlers {
	return &currencyHandlers{
		cfg: cfg,
		currencyUseCase: currencyUseCase,
		logger: logger,
	}
}
