package v1

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Hamedblue1381/restaurant-reserve/models"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

var userHandler *models.UserHandler

type ErrorResponse struct {
	Error string `json:"error" example:"Description of the error occurred"`
}

func InitializedUserHandler(db *gorm.DB) {
	userHandler = models.NewUserHandler(db)
}

// @Summary Get a Single User
// @Description Retrieves details of a single user by their unique identifier.
// @Tags user
// @Produce json
// @Param id path int true "User ID" Format(int64)
// @Security Bearer
// @Success 200 {object} models.User "The details of the user including ID, name, email, telephone, and role."
// @Failure 400 {object} ErrorResponse "Invalid user ID format."
// @Failure 404 {object} ErrorResponse "User not found with the specified ID."
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	idUint := uint(idInt)

	user, err := userHandler.GetUser(idUint)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Get All Users
// @Description Retrieves a list of all users in the system.
// @Tags user
// @Produce json
// @Security Bearer
// @Success 200 {array} models.User "An array of user objects."
// @Failure 500 {object} ErrorResponse "Internal server error while fetching users."
// @Router /users [get]
func GetUsers(c *gin.Context) {
	users, err := userHandler.GetUsers()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users!"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Create a New User
// @Description Adds a new user to the system with the provided details.
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.User true "User Registration Details"
// @Security Bearer
// @Success 201 {object} models.User "The created user's details, including their unique identifier."
// @Failure 400 {object} ErrorResponse "Invalid input format for user details."
// @Failure 500 {object} ErrorResponse "Internal server error while creating the user."
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := userHandler.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user!"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// @Summary Update a User
// @Description Updates the details of an existing user identified by their ID.
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID" Format(int64)
// @Param user body models.User true "Updated User Details"
// @Security Bearer
// @Success 200 {object} models.User "The updated user's details."
// @Failure 400 {object} ErrorResponse "Invalid input format for user details or invalid user ID."
// @Failure 500 {object} ErrorResponse "Internal server error while updating the user."
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}
	idUint := uint(idInt)

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = userHandler.UpdateUser(idUint, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Delete a User
// @Description Removes a user from the system by their unique identifier.
// @Tags user
// @Produce json
// @Param id path int true "User ID" Format(int64)
// @Security Bearer
// @Success 204 "User successfully deleted, no content to return."
// @Failure 400 {object} ErrorResponse "Invalid user ID format."
// @Failure 500 {object} ErrorResponse "Internal server error while deleting the user."
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}
	idUint := uint(idInt)

	err = userHandler.DeleteUser(idUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user"})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Get my profile
// @Description Retrieves the details of the currently authenticated user.
// @Tags user
// @Produce json
// @Security Bearer
// @Success 200 {object} models.User "The details of the currently authenticated user."
// @Failure 404 {object} ErrorResponse "User not found."
// @Router /me [get]
func GetMe(c *gin.Context) {
	userId, _ := c.Get("id")
	if userId == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User must be logged in to view reservations"})
		return
	}
	user, err := userHandler.GetUser(userId.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary Get my profile QR CODE
// @Description Retrieves the QR Code of the currently authenticated user.
// @Tags user
// @Produce json
// @Security Bearer
// @Success 200 {object} models.User "The QR CODE of the currently authenticated user."
// @Failure 404 {object} ErrorResponse "User not found."
// @Router /me/qr [get]
func GetMeQR(c *gin.Context) {
	userId, _ := c.Get("id")
	if userId == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User must be logged in to view reservations"})
		return
	}
	urlToUser := fmt.Sprintf("http://%s/api/v1/users/%v", os.Getenv("BASE_URL"), userId)

	// Generate QR code
	qrCode, err := qrcode.Encode(urlToUser, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	// Return QR code as response
	c.Data(http.StatusOK, "image/png", qrCode)
}
