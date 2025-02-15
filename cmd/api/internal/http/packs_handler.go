package http

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"re_partners_software_coding_challenge_go/cmd/api/internal/domain/orderpacks"
)

// PacksHandler handles requests related to order packs calculations.
type PacksHandler struct {
	OrderPacksCalculator orderpacks.UseCaseCalculateOrderPacks
}

// ShowOrderPacksCalculator displays the order packs calculator form.
func (h *PacksHandler) ShowOrderPacksCalculator(c *gin.Context) {
	data := map[string]interface{}{
		"Error":      c.Query("error"), // Error message from query parameter.
		"packSizes":  1,                // Default pack sizes (for display).
		"orderItems": 1,                // Default order items (for display).
	}

	// Render the HTML template.
	html, err := renderTemplate(data)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// renderTemplate renders the HTML template with the provided data.
func renderTemplate(pageData interface{}) (string, error) {
	// Read and parse the HTML template file.
	tmpl, err := template.ParseFiles("./static/order_packs_calculator_form.html")
	if err != nil {
		return "", fmt.Errorf("Error parsing template: %v ", err)
	}

	var renderedTemplate strings.Builder

	err = tmpl.Execute(&renderedTemplate, pageData)
	if err != nil {
		return "", fmt.Errorf("Error parsing template: %v ", err)
	}

	resultString := renderedTemplate.String()

	return resultString, nil
}

// CalculateOrderPacks calculates order packs and displays the results.
func (h *PacksHandler) CalculateOrderPacks(c *gin.Context) {
	// Get form data.
	orderPackForm, err := getOrderPacksForm(c)
	if err != nil {
		c.Redirect(http.StatusFound, "/?error="+err.Error()) // Redirect with error.
		return
	}

	// Parse pack sizes.
	packSizesStr := strings.Split(orderPackForm.PackSizes, ",")
	packSizes := make([]uint64, len(packSizesStr))
	for i, psStr := range packSizesStr {
		psStr = strings.TrimSpace(psStr)

		psUint, err := strconv.ParseUint(psStr, 10, 64)
		if err != nil {
			c.Redirect(http.StatusFound, "/?error=invalid pack sizes") // Redirect with error.
			return
		}

		packSizes[i] = psUint
	}

	// Parse order items.
	orderItems, err := strconv.ParseUint(orderPackForm.OrderItems, 10, 64)
	if err != nil {
		c.Redirect(http.StatusFound, "/?error=invalid order items") // Redirect with error.
		return
	}

	// Calculate order packs.
	packs, err := h.OrderPacksCalculator.Run(orderpacks.UseCaseCalculateOrderPacksRequest{
		PackSizes:  packSizes,
		OrderItems: orderItems,
	})
	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/?error=%s", err)) // Redirect with error.
		return
	}

	// Prepare template data.
	data := map[string]interface{}{
		"Error":      c.Query("error"),
		"packSizes":  orderPackForm.PackSizes,
		"orderItems": orderPackForm.OrderItems,
	}

	var results []map[string]uint64

	for key, value := range packs {
		result := map[string]uint64{
			"Pack":     key,
			"Quantity": value,
		}

		results = append(results, result)
	}

	data["Results"] = results

	// Render the HTML template.
	html, err := renderTemplate(data)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// OrderPacks represents the form data for order packs.
type OrderPacks struct {
	PackSizes  string `form:"packSizes"   binding:"required"` // Pack sizes (comma-separated).
	OrderItems string `form:"orderItems" binding:"required"`  // Number of items to order.
}

// getOrderPacksForm extracts order pack data from the request body
func getOrderPacksForm(c *gin.Context) (OrderPacks, error) {
	if c.Request.Body == nil {
		return OrderPacks{}, fmt.Errorf("body cannot be nil")
	}

	form := OrderPacks{}

	// Parse the form data from the request body and populate the OrderPacks struct.
	if err := binding.FormPost.Bind(c.Request, &form); err != nil {
		log.Println(err.Error())
		return OrderPacks{}, err
	}

	return form, nil
}
