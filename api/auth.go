package api

import (
	"net/http"

	"github.com/Hamedblue1381/restaurant-reserve/middleware"
	"github.com/Hamedblue1381/restaurant-reserve/models"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

var userHandler *models.UserHandler

func InitializedAuthHandler(db *gorm.DB) {
	userHandler = models.NewUserHandler(db)
}

type RegisterDetails struct {
	Name      string `json:"name" example:"John Doe"`
	Telephone string `json:"telephone" example:"123-456-7890"`
	Email     string `json:"email" example:"john.doe@example.com"`
	Password  string `json:"password" example:"securePassword123"`
	Role      string `json:"role" example:"user"`
}

type RegisterResponse struct {
	Message string `json:"message" example:"User registered successfully"`
}

// @Summary Register a new user
// @Description Creates a new user account with the provided details. Upon successful creation, the user can log in with their credentials.
// @Tags authentication
// @Accept json
// @Produce json
// @Param user body RegisterDetails true "Register Credentials"
// @Success 200 {object} RegisterResponse "Confirmation of successful registration."
// @Failure 400 {object} ErrorResponse "The request was formatted incorrectly or missing required fields."
// @Failure 500 {object} ErrorResponse "Internal server error, unable to process the request."
// @Router /auth/register [post]
func Register(c *gin.Context) {

	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input format! please check the input format"})
		return
	}

	err := userHandler.CreateUser(&newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	token, err := middleware.GenerateToken(newUser.Email, newUser.ID, newUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "User registered successfully",
			"token":   token,
		},
	)
}

type LoginDetails struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

type LoginResponse struct {
	Token   string `json:"token" example:""`
	Message string `json:"message" example:"Login successful"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Error message"`
}

// Login a user
// @Summary User Login
// @Description Authenticates a user by their email and password, returning a JWT token for authorized access to protected endpoints if successful.
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body LoginDetails true "Login Credentials"
// @Success 200 {object} LoginResponse "An object containing a JWT token for authentication and a message indicating successful login."
// @Failure 400 {object} ErrorResponse "The request was formatted incorrectly or missing required fields."
// @Failure 401 {object} ErrorResponse "Authentication failed due to invalid login credentials."
// @Failure 404 {object} ErrorResponse "The specified user was not found in the system."
// @Failure 500 {object} ErrorResponse "Internal server error, unable to process the request."
// @Router /auth/signin [post]
func Login(c *gin.Context) {

	var loginDetails LoginDetails
	if err := c.ShouldBindJSON(&loginDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input format! please check the input format"})
		return
	}

	user, err := userHandler.GetUserByEmail(loginDetails.Email)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Authentication failed"})
		return
	}

	if !userHandler.CheckPassword(user.Email, loginDetails.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password is incorrect!"})
		return
	}

	token, err := middleware.GenerateToken(user.Email, user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"token":   token,
			"message": "Login successful",
		},
	)
}
