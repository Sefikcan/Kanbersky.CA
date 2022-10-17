package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/sefikcan/kanbersky.ca/docs"
	"github.com/sefikcan/kanbersky.ca/internal/currency/handlers"
	"github.com/sefikcan/kanbersky.ca/internal/currency/repository"
	"github.com/sefikcan/kanbersky.ca/internal/currency/usecase"
	mw "github.com/sefikcan/kanbersky.ca/internal/middleware"
	"github.com/sefikcan/kanbersky.ca/pkg/metric"
	"github.com/sefikcan/kanbersky.ca/pkg/util"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	metrics, err := metric.CreateMetrics(s.cfg.Metric.Url, s.cfg.Metric.ServiceName)
	if err != nil {
		s.logger.Errorf("CreateMetrics error: %s", err)
	}
	s.logger.Info("Metrics available URL: %s, ServiceName: %s", s.cfg.Metric.Url, s.cfg.Metric.ServiceName)

	currencyRepository := repository.NewCurrencyRepository(s.db)
	currencyRedisRepository := repository.NewCurrencyRedisRepository(s.redisClient)

	currencyUseCase := usecase.NewCurrencyUseCase(s.cfg, currencyRepository, currencyRedisRepository, s.logger)

	currencyHandler := handlers.NewCurrencyHandler(s.cfg, currencyUseCase, s.logger)

	middlewareManager := mw.NewMiddlewareManager(s.cfg, s.logger)
	e.Use(middlewareManager.RequestLoggerMiddleware)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, //1kb
		DisablePrintStack: true,
		DisableStackAll: true,
	}))
	e.Use(middleware.RequestID())
	e.Use(middlewareManager.MetricsMiddleware(metrics))
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := e.Group("/api/v1")
	health := v1.Group("/health")
	currencyGroup := v1.Group("/currencies")

	handlers.MapCurrencyRoutes(currencyGroup, currencyHandler)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check RequestID: %s", util.GetRequestId(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}
