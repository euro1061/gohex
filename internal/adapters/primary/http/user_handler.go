package http

import (
	"os"
	"strconv"
	"time"

	"github.com/euro1061/gohex/internal/domain/entity"
	"github.com/euro1061/gohex/internal/ports/input"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	userService input.UserService
}

func NewUserHandler(service input.UserService) *UserHandler {
	return &UserHandler{
		userService: service,
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	user := new(entity.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.userService.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.userService.GetUser(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user := new(entity.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.userService.UpdateUser(user, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.userService.DeleteUser(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	credentials := new(entity.Credentials)
	if err := c.BodyParser(credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.userService.Login(credentials.Email, credentials.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	claims := jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
	})
}

// Implement other handler methods...
