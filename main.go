package main

import (
	"log"
	"os"
	"strconv"

	"GO-X1/auth"
	"GO-X1/connectDB"
	"GO-X1/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// Global variables
var (
	db       *gorm.DB
	validate *validator.Validate
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize validator
	validate = validator.New()

	// Connect to database
	var err error
	db, err = connectdb.ConnectDB()
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		log.Println("Continuing without database connection...")
	} else {
		log.Println("Successfully connected to database!")
		
		// Auto migrate the schema
		if err := db.AutoMigrate(&models.User{}); err != nil {
			log.Printf("Warning: Failed to migrate database: %v", err)
		} else {
			log.Println("Database migration completed!")
		}
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return ctx.Status(code).JSON(APIResponse{
				Success: false,
				Message: "An error occurred",
				Error:   err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Routes
	setupRoutes(app)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func setupRoutes(app *fiber.App) {
	// Health check
	app.Get("/", healthCheck)
	app.Get("/health", healthCheck)
	
	// Utility endpoints
	app.Get("/uuid", generateUUID)

	// API group
	api := app.Group("/api/v1")
	
	// Auth routes
	authRoutes := api.Group("/auth")
	authRoutes.Post("/login", loginUser)

	// User routes
	users := api.Group("/users")
	users.Post("/", createUser)

	// Protected user routes
	protectedUsers := users.Use(auth.AuthMiddleware)
	protectedUsers.Get("/", getUsers)
	protectedUsers.Get("/:id", getUserByID)
	protectedUsers.Put("/:id", updateUser)
	protectedUsers.Delete("/:id", deleteUser)
}

// Login user
func loginUser(c *fiber.Ctx) error {
	if db == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(APIResponse{
			Success: false,
			Message: "Database not available",
			Error:   "Database connection not established",
		})
	}

	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	// Find user by email
	var user models.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(APIResponse{
				Success: false,
				Message: "Invalid credentials",
				Error:   "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Success: false,
			Message: "Failed to fetch user",
			Error:   err.Error(),
		})
	}

	// Check password
	if !auth.CheckPasswordHash(req.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(APIResponse{
			Success: false,
			Message: "Invalid credentials",
			Error:   "Incorrect password",
		})
	}

	// Generate JWT
	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Success: false,
			Message: "Failed to generate token",
			Error:   err.Error(),
		})
	}

	return c.JSON(APIResponse{
		Success: true,
		Message: "Login successful",
		Data: models.LoginResponse{
			Token: token,
			User:  user.ToResponse(),
		},
	})
}

// Health check endpoint
func healthCheck(c *fiber.Ctx) error {
	return c.JSON(APIResponse{
		Success: true,
		Message: "API is running successfully",
		Data: map[string]interface{}{
			"status": "healthy",
			"service": "GO-X1 REST API",
			"database": db != nil,
		},
	})
}

// Generate UUID endpoint
func generateUUID(c *fiber.Ctx) error {
	uuid := auth.GenUuid()
	return c.JSON(APIResponse{
		Success: true,
		Message: "UUID generated successfully",
		Data:    map[string]string{"uuid": uuid},
	})
}

// Get all users
func getUsers(c *fiber.Ctx) error {
	if db == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(APIResponse{
			Success: false,
			Message: "Database not available",
			Error:   "Database connection not established",
		})
	}

	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Success: false,
			Message: "Failed to fetch users",
			Error:   err.Error(),
		})
	}

	// Convert to response format
	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	return c.JSON(APIResponse{
		Success: true,
		Message: "Users fetched successfully",
		Data:    userResponses,
	})
}

// Create a new user
func createUser(c *fiber.Ctx) error {
	if db == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(APIResponse{
			Success: false,
			Message: "Database not available",
			Error:   "Database connection not established",
		})
	}

	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Success: false,
			Message: "Failed to hash password",
			Error:   err.Error(),
		})
	}

	// Create user
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Success: false,
			Message: "Failed to create user",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Success: true,
		Message: "User created successfully",
		Data:    user.ToResponse(),
	})
}

// Get user by ID
func getUserByID(c *fiber.Ctx) error {
	if db == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(APIResponse{
			Success: false,
			Message: "Database not available",
			Error:   "Database connection not established",
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   err.Error(),
		})
	}

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(APIResponse{
				Success: false,
				Message: "User not found",
				Error:   "User with the specified ID does not exist",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Success: false,
			Message: "Failed to fetch user",
			Error:   err.Error(),
		})
	}

	return c.JSON(APIResponse{
		Success: true,
		Message: "User fetched successfully",
		Data:    user.ToResponse(),
	})
}

// Update user
func updateUser(c *fiber.Ctx) error {
	if db == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(APIResponse{
			Success: false,
			Message: "Database not available",
			Error:   "Database connection not established",
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   err.Error(),
		})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	// Find user
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(APIResponse{
				Success: false,
				Message: "User not found",
				Error:   "User with the specified ID does not exist",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Success: false,
			Message: "Failed to fetch user",
			Error:   err.Error(),
		})
	}

	// Update fields
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		user.Password = req.Password // In production, hash the password
	}

	if err := db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Success: false,
			Message: "Failed to update user",
			Error:   err.Error(),
		})
	}

	return c.JSON(APIResponse{
		Success: true,
		Message: "User updated successfully",
		Data:    user.ToResponse(),
	})
}

// Delete user
func deleteUser(c *fiber.Ctx) error {
	if db == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(APIResponse{
			Success: false,
			Message: "Database not available",
			Error:   "Database connection not established",
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   err.Error(),
		})
	}

	// Check if user exists
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(APIResponse{
				Success: false,
				Message: "User not found",
				Error:   "User with the specified ID does not exist",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Success: false,
			Message: "Failed to fetch user",
			Error:   err.Error(),
		})
	}

	if err := db.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Success: false,
			Message: "Failed to delete user",
			Error:   err.Error(),
		})
	}

	return c.JSON(APIResponse{
		Success: true,
		Message: "User deleted successfully",
		Data:    map[string]interface{}{"deleted_user_id": id},
	})
}
