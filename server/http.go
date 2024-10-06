package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func registerHandlers(e *echo.Echo, h *handler) {
	g := e.Group("/ton-market")
	g.POST("/generate-payload", h.PayloadHandler, middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.POST},
	}))
	g.POST("/check-proof", h.ProofHandler, middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.POST},
	}))
	g.Use(middleware.CORS())
	g.GET("/get-address", h.GetAccountInfo, middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET},
	}), middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper:    middleware.DefaultSkipper,
		KeyLookup:  "cookie:AuthToken",
		AuthScheme: "Bearer",
		Validator:  h.validateUser,
	}))

	g.GET("/get-events", h.GetEvents, middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET},
	}))
}
