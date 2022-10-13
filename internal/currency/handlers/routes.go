package handlers

import "github.com/labstack/echo/v4"

func MapCurrencyRoutes(currencyRouteGroup *echo.Group, c CurrencyHandlers) {
	currencyRouteGroup.POST("/", c.Create())
	currencyRouteGroup.PUT("/:id", c.Update())
	currencyRouteGroup.DELETE("/:id", c.Delete())
	currencyRouteGroup.GET("/:id", c.GetById())
	currencyRouteGroup.GET("", c.GetAll())
}
