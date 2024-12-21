package http

import (
	"time"

	"github.com/euro1061/gohex/internal/application"
	"github.com/euro1061/gohex/internal/dto"
	"github.com/euro1061/gohex/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service   *application.UserService
	validator *validator.Validate
}

func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	app.Post("/register", h.Register)
	app.Post("/login", h.Login)
	app.Put("/users/profile", middleware.Auth(), h.UpdateProfile)
	app.Get("/users/profile", middleware.Auth(), h.GetProfile)
	app.Post("/logout", middleware.Auth(), h.Logout)
}

// @Summary Register a new user
// @Description Register a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.UserRegisterRequest true "User registration info"
// @Success 201 {object} Response{data=dto.UserResponse}
// @Failure 400 {object} ErrorResponse
// @Router /register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req dto.UserRegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	// Convert DTO to domain model and register
	user := req.ToUser()
	if err := h.service.Register(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Message: "Failed to register user",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: "User registered successfully",
		Data:    dto.UserResponseFromUser(user),
	})
}

// @Summary Login user
// @Description Authenticate a user and return a JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.UserLoginRequest true "User credentials"
// @Success 200 {object} Response{data=string}
// @Failure 401 {object} ErrorResponse
// @Router /login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.UserLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	token, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Success: false,
			Message: "Login failed",
			Error:   err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
		SameSite: "lax",
	})

	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: "Login successful",
		Data:    token,
	})
}

// @Summary Update user profile
// @Description Update the authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body dto.UserUpdateRequest true "User update info"
// @Success 200 {object} Response{data=dto.UserResponse}
// @Failure 401 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	var req dto.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Message: "Invalid request format",
			Error:   err.Error(),
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	// Get user from token
	token := c.Locals("token").(string)
	user, err := h.service.GetUserFromToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Success: false,
			Message: "Invalid or expired token",
			Error:   err.Error(),
		})
	}

	// Update user fields
	user.Name = req.Name
	user.Username = req.Username
	user.Gender = req.Gender
	user.Email = req.Email

	// Update user
	if err := h.service.Update(user); err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "username already taken" || err.Error() == "email already taken" {
			status = fiber.StatusConflict
		}
		return c.Status(status).JSON(ErrorResponse{
			Success: false,
			Message: "Failed to update profile",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: "Profile updated successfully",
		Data:    dto.UserResponseFromUser(user),
	})
}

// @Summary Get user profile
// @Description Get the authenticated user's profile information
// @Tags users
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} Response{data=dto.UserResponse}
// @Failure 401 {object} ErrorResponse
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	// Get user from token
	token := c.Locals("token").(string)
	user, err := h.service.GetUserFromToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
			Success: false,
			Message: "Invalid or expired token",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    dto.UserResponseFromUser(user),
	})
}

// @Summary Logout user
// @Description Log out the authenticated user
// @Tags users
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} Response
// @Failure 401 {object} ErrorResponse
// @Router /logout [post]
func (h *UserHandler) Logout(c *fiber.Ctx) error {
	c.ClearCookie("token")
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: "Logout successful",
	})
}
