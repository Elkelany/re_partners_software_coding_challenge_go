package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"re_partners_software_coding_challenge_go/cmd/api/internal/domain/orderpacks"
)

// NewServer creates and configures a new HTTP server.
func NewServer() *http.Server {
	var handler = PacksHandler{
		OrderPacksCalculator: orderpacks.UseCaseCalculateOrderPacks{},
	}

	ginEngine := gin.Default()

	// - GET "/" serves the order packs calculator form.
	ginEngine.GET("/", handler.ShowOrderPacksCalculator)

	// - POST "/calculate-order-packs" handles the order packs calculation request.
	ginEngine.POST("/calculate-order-packs", handler.CalculateOrderPacks)

	return &http.Server{
		Addr:    ":8088", // Listen on port 8088
		Handler: ginEngine,
	}
}
