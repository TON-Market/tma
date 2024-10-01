package main

import (
	"github.com/TON-Market/tma/server/config"
	echojwt "github.com/labstack/echo-jwt"
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
	}), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.Config.Proof.PayloadSignatureKey),
	}))
}
