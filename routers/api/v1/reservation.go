package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Hamedblue1381/restaurant-reserve/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SuccessResponse struct {
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
}

var reservationHandler *models.ReservationHandler

func InitializeReservationHandler(db *gorm.DB) {
	reservationHandler = models.NewReservationHandler(db)
}

// @Summary Create a reservation
// @Description Create a new reservation
// @Tags reservation
// @Accept json
// @Produce json
// @Param reservation body models.Reservation true "Reservation details"
// @Security Bearer
// @Success 200 {object} SuccessResponse "The created reservation's date"
// @Failure 403 {object} ErrorResponse "User must be logged in to update a reservation"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /reservation [post]
func CreateReservation(c *gin.Context) {
	var reservation models.Reservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Check if user is logged in
	userId, _ := c.Get("id")
	if userId == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User must be logged in to make a reservation"})
		return
	}

	// Convert user ID to uint
	userIdUint, _ := userId.(uint)
	// Check if user has more than 3 unpaid reservations

	if err := reservationHandler.IsBlackListed(userIdUint); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err})
		c.Abort()
		return
	}

	// Set user ID for the reservation
	reservation.UserID = userIdUint

	if err := reservationHandler.Reserve(&reservation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reservation"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Reservation created successfully",
		Date:    reservation.Date,
	})
}

// @Summary Delete a reservation
// @Description Delete a reservation by ID
// @Tags reservation
// @Produce json
// @Param id path int true "Reservation ID"
// @Security Bearer
// @Success 204 "No content"
// @Failure 400 {object} ErrorResponse "Invalid reservation ID format"
// @Failure 403 {object} ErrorResponse "User must be logged in to update a reservation"
// @Failure 404 {object} ErrorResponse "Reservation not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /reservation/{id} [delete]
func DeleteReservation(c *gin.Context) {
	// Authenticate user
	userId, _ := c.Get("id")
	if userId == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User must be logged in to delete a reservation"})
		return
	}

	// Parse reservation ID
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation ID"})
		return
	}
	idUint := uint(idInt)

	// Delete reservation
	if err := reservationHandler.DeleteReservation(idUint); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Update a reservation
// @Description Update an existing reservation by ID
// @Tags reservation
// @Accept json
// @Produce json
// @Param id path int true "Reservation ID"
// @Param reservation body models.Reservation true "Reservation details"
// @Security Bearer
// @Success 200 {object} models.Reservation "The updated reservation"
// @Failure 403 {object} ErrorResponse "User must be logged in to update a reservation"
// @Failure 400 {object} ErrorResponse "Invalid request format"
// @Failure 404 {object} ErrorResponse "Reservation not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /reservation/{id} [put]
func UpdateReservation(c *gin.Context) {
	// Authenticate user
	userId, _ := c.Get("id")
	if userId == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User must be logged in to update a reservation"})
		return
	}

	// Parse reservation ID
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation ID"})
		return
	}
	idUint := uint(idInt)

	// Parse updated reservation data
	var updatedReservation models.Reservation
	if err := c.ShouldBindJSON(&updatedReservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Update reservation
	if err := reservationHandler.UpdateReservation(idUint, &updatedReservation); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	c.JSON(http.StatusOK, updatedReservation)
}

// @Summary get reservations
// @Description List reservations based on provided start and end dates
// @Tags reservation
// @Produce json
// @Param start_date query string false "Start date (format: yyyy-mm-dd)"
// @Param end_date query string false "End date (format: yyyy-mm-dd)"
// @Security Bearer
// @Success 200 {array} models.Reservation "List of reservations"
// @Failure 400 {object} ErrorResponse "Invalid date format"
// @Failure 403 {object} ErrorResponse "User must be logged in to update a reservation"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /reservations [get]
func GetReservations(c *gin.Context) {
	// Authenticate user
	userId, _ := c.Get("id")
	if userId == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User must be logged in to view reservations"})
		return
	}

	// Parse start and end dates
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Parse start date
	var startDate time.Time
	if startDateStr != "" {
		var err error
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
			return
		}
	}

	// Parse end date
	var endDate time.Time
	if endDateStr != "" {
		var err error
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}
	}

	// List reservations
	reservations, err := reservationHandler.ListReservations(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list reservations"})
		return
	}

	c.JSON(http.StatusOK, reservations)
}

// @Summary Get a single reservation
// @Description Retrieve details of a single reservation by its unique identifier
// @Tags reservation
// @Produce json
// @Param id path int true "Reservation ID"
// @Security Bearer
// @Success 200 {object} models.Reservation "The reservation details"
// @Failure 400 {object} ErrorResponse "Invalid reservation ID format"
// @Failure 403 {object} ErrorResponse "User must be logged in to update a reservation"
// @Failure 404 {object} ErrorResponse "Reservation not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /reservation/{id} [get]
func GetReservation(c *gin.Context) {
	// Authenticate user
	userId, _ := c.Get("id")
	if userId == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User must be logged in to view a reservation"})
		return
	}

	// Parse reservation ID
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation ID"})
		return
	}
	idUint := uint(idInt)

	// Get reservation
	reservation, err := reservationHandler.GetReservation(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

// GetUserReservations retrieves all reservations for a given user ID.
// @Summary Get User's Reservations
// @Description Retrieves a list of reservations associated with a specific user.
// @Tags reservation
// @Produce json
// @Param userId path int true "User ID"
// @Security Bearer
// @Success 200 {array} models.Reservation "An array of reservation objects for the user."
// @Failure 400 {object} ErrorResponse "Invalid user ID format."
// @Failure 404 {object} ErrorResponse "Reservations not found for the specified user ID."
// @Router /users/{userId}/reservations [get]
func GetUserReservations(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	uid, err := strconv.ParseUint(userID, 10, 32) // Convert userID from string to uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing user ID"})
		return
	}

	reservations, err := reservationHandler.GetReservationsByUserID(uint(uid)) // Correctly cast to uint now
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching reservations for user"})
		return
	}

	c.JSON(http.StatusOK, reservations)
}
