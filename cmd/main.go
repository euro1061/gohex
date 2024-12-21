package main

import (
	"fmt"
	"log"
	"os"

	"github.com/euro1061/gohex/internal/adapters/repository/postgres"
	"github.com/euro1061/gohex/internal/application"
	"github.com/euro1061/gohex/internal/ports/http"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	postgresDB "gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/euro1061/gohex/docs" // Import generated docs
)

// @title GoHex API
// @version 1.0
// @description This is a sample API server for GoHex application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgresDB.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Initialize repositories
	productRepo := postgres.NewProductRepository(db)
	userRepo := postgres.NewUserRepository(db)

	// Initialize services
	productService := application.NewProductService(productRepo)
	userService := application.NewUserService(userRepo)

	// Initialize HTTP handlers
	productHandler := http.NewProductHandler(productService)
	userHandler := http.NewUserHandler(userService)

	// Setup Fiber app
	app := fiber.New()

	// CORS middleware with more secure configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000", // เปลี่ยนเป็น domain ของ frontend จริงๆ
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	}))

	// Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Register routes
	productHandler.RegisterRoutes(app)
	userHandler.RegisterRoutes(app)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func initDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	return gorm.Open(postgresDB.Open(dsn), &gorm.Config{})
}
