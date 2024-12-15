package main

import (
	"log"

	"github.com/euro1061/gohex/configs"
	"github.com/euro1061/gohex/internal/adapters/primary/http"
	"github.com/euro1061/gohex/internal/adapters/primary/http/middleware"
	"github.com/euro1061/gohex/internal/adapters/secondary/postgres"
	"github.com/euro1061/gohex/internal/domain/entity"
	"github.com/euro1061/gohex/internal/ports/input"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	pgDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := configs.LoadConfig()

	// Initialize Database
	db, err := gorm.Open(pgDriver.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Auto Migrate
	db.AutoMigrate(&entity.User{})

	authMiddleware := middleware.NewAuthMiddleware()

	userRepo := postgres.NewUserRepository(db)
	userService := input.NewUserService(userRepo)
	userHandler := http.NewUserHandler(userService)

	app := fiber.New()

	// Setup Routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/users/reigster", userHandler.CreateUser)
	v1.Post("/users/login", userHandler.Login)

	users := v1.Group("/users")
	users.Use(authMiddleware.VerifyToken)
	// users.Post("/", userHandler.CreateUser)
	users.Get("/:id", userHandler.GetUser)
	users.Get("/", userHandler.GetAllUsers)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
	// users.Post("/login", userHandler.Login)

	// Start Server
	log.Fatal(app.Listen(config.ServerPort))
}
