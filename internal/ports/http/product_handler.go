package http

import (
	"encoding/json"
	"strconv"

	"github.com/euro1061/gohex/internal/application"
	"github.com/euro1061/gohex/internal/domain"
	"github.com/euro1061/gohex/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type ProductHandler struct {
	service *application.ProductService
}

func NewProductHandler(service *application.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (h *ProductHandler) RegisterRoutes(app *fiber.App) {
	app.Post("/products", middleware.Auth(), h.CreateProduct)
	app.Get("/products", h.GetAllProducts)
	app.Get("/products/:id", h.GetProduct)
	app.Put("/products/:id", middleware.Auth(), h.UpdateProduct)
	app.Delete("/products/:id", middleware.Auth(), h.DeleteProduct)
}

// @Summary Create a new product
// @Description Create a new product with the provided information
// @Tags products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param product body domain.Product true "Product info"
// @Success 201 {object} Response{data=domain.Product}
// @Failure 401 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var product domain.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(ErrorResponse{
			Success: false,
			Message: "Failed to create product",
			Error:   "Invalid request payload",
		})
	}

	createdProduct, err := h.service.CreateProduct(product.Name, product.Description, product.Price)
	if err != nil {
		return c.Status(500).JSON(ErrorResponse{
			Success: false,
			Message: "Failed to create product",
			Error:   err.Error(),
		})
	}

	return c.Status(201).JSON(Response{
		Success: true,
		Message: "Product created successfully",
		Data:    createdProduct,
	})
}

// @Summary Get all products
// @Description Get a list of all products
// @Tags products
// @Produce json
// @Success 200 {object} Response{data=[]domain.Product}
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.service.GetAllProducts()
	if err != nil {
		return c.Status(500).JSON(ErrorResponse{
			Success: false,
			Message: "Failed to get products",
			Error:   err.Error(),
		})
	}

	return c.JSON(Response{
		Success: true,
		Message: "Products retrieved successfully",
		Data:    products,
	})
}

// @Summary Get a product
// @Description Get a product by its ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} Response{data=domain.Product}
// @Failure 404 {object} ErrorResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(ErrorResponse{
			Success: false,
			Message: "Failed to get product",
			Error:   "Invalid product ID",
		})
	}

	product, err := h.service.GetProduct(uint(id))
	if err != nil {
		return c.Status(404).JSON(ErrorResponse{
			Success: false,
			Message: "Failed to get product",
			Error:   "Product not found",
		})
	}

	return c.JSON(Response{
		Success: true,
		Message: "Product retrieved successfully",
		Data:    product,
	})
}

// @Summary Update a product
// @Description Update a product with the provided information
// @Tags products
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Product ID"
// @Param product body domain.Product true "Product info"
// @Success 200 {object} Response{data=domain.Product}
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(ErrorResponse{
			Success: false,
			Message: "Failed to update product",
			Error:   "Invalid product ID",
		})
	}

	var product domain.Product
	if err := json.Unmarshal(c.Body(), &product); err != nil {
		return c.Status(400).JSON(ErrorResponse{
			Success: false,
			Message: "Failed to update product",
			Error:   "Invalid request payload",
		})
	}

	product.ID = uint(id)
	if err := h.service.UpdateProduct(&product); err != nil {
		return c.Status(500).JSON(ErrorResponse{
			Success: false,
			Message: "Failed to update product",
			Error:   err.Error(),
		})
	}

	return c.JSON(Response{
		Success: true,
		Message: "Product updated successfully",
		Data:    product,
	})
}

// @Summary Delete a product
// @Description Delete a product by its ID
// @Tags products
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Product ID"
// @Success 200 {object} Response
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(ErrorResponse{
			Success: false,
			Message: "Failed to delete product",
			Error:   "Invalid product ID",
		})
	}

	if err := h.service.DeleteProduct(uint(id)); err != nil {
		return c.Status(500).JSON(ErrorResponse{
			Success: false,
			Message: "Failed to delete product",
			Error:   err.Error(),
		})
	}

	return c.JSON(Response{
		Success: true,
		Message: "Product deleted successfully",
		Data:    nil,
	})
}
