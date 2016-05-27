package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"

	"github.com/nvhbk16k53/charging/handler"
)

func main() {
	// Create new echo object
	server := echo.New()
	// Enable debuging
	server.SetDebug(true)

	// Add API endpoints
	server.Get("/accounts/:id/charges", handler.ListChargeByAccountID)
	server.Post("/charges", handler.CreateChargeRequest)
	server.Put("/charges/:id/approve", handler.UpdateChargeRequest)

	// Run server
	server.Run(standard.New(":9909"))
}
