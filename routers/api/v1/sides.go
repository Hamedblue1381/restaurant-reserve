package v1

import (
	"net/http"
	"strconv"

	"github.com/Hamedblue1381/restaurant-reserve/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var sidesHandler *models.SidesHandler

func InitializedSidesHandler(db *gorm.DB) {
	sidesHandler = models.NewSidesHandler(db)
}

// @Summary Get a Single Side Dish
// @Description Retrieves details of a single side dish by their unique identifier.
// @Tags sides
// @Produce json
// @Param id path int true "Sides ID" Format(int64)
// @Security Bearer
// @Success 200 {object} models.Sides "The details of the sides including ID, name, quantity."
// @Failure 400 {object} ErrorResponse "Invalid sides ID format."
// @Failure 404 {object} ErrorResponse "Sides not found with the specified ID."
// @Router /sides/{id} [get]
func GetSide(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sides id"})
		return
	}

	idUint := uint(idInt)

	user, err := sidesHandler.GetSide(idUint)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Side Dish not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Get All Sides
// @Description Retrieves a list of all side dishes in the system.
// @Tags sides
// @Produce json
// @Security Bearer
// @Success 200 {array} models.Sides "An array of sides objects."
// @Failure 500 {object} ErrorResponse "Internal server error while fetching sides."
// @Router /sides [get]
func GetSides(c *gin.Context) {
	sides, err := sidesHandler.GetSides()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching sides!"})
		return
	}

	c.JSON(http.StatusOK, sides)
}

// @Summary Create a New Sides
// @Description Adds a new side dish to the system with the provided details.
// @Tags sides
// @Accept json
// @Produce json
// @Param sides body models.Sides true "Sides Details"
// @Security Bearer
// @Success 201 {object} models.Sides "The created Side's details, including their unique identifier."
// @Failure 400 {object} ErrorResponse "Invalid input format for Sides."
// @Failure 500 {object} ErrorResponse "Internal server error while creating the sides."
// @Router /sides [post]
func CreateSides(c *gin.Context) {
	var side models.Sides

	if err := c.ShouldBindJSON(&side); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := sidesHandler.CreateSides(&side); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating side dish!"})
		return
	}

	c.JSON(http.StatusCreated, side)
}

// @Summary Update a Side Dish
// @Description Updates the details of an existing side dish identified by their ID.
// @Tags sides
// @Accept json
// @Produce json
// @Param id path int true "Side ID" Format(int64)
// @Param sides body models.Sides true "Updated Sides Details"
// @Security Bearer
// @Success 200 {object} models.Sides "The updated side's details."
// @Failure 400 {object} ErrorResponse "Invalid input format for user details or invalid sides ID."
// @Failure 500 {object} ErrorResponse "Internal server error while updating the sides."
// @Router /sides/{id} [put]
func UpdateSides(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid side id"})
		return
	}
	idUint := uint(idInt)

	var side models.Sides

	if err := c.ShouldBindJSON(&side); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = sidesHandler.UpdateSides(idUint, &side)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating side dish"})
		return
	}

	c.JSON(http.StatusOK, side)
}

// @Summary Delete a Side Dish
// @Description Removes a side dish from the system by their unique identifier.
// @Tags sides
// @Produce json
// @Param id path int true "Sides ID" Format(int64)
// @Security Bearer
// @Success 204 "Side Dish successfully deleted, no content to return."
// @Failure 400 {object} ErrorResponse "Invalid sides ID format."
// @Failure 500 {object} ErrorResponse "Internal server error while deleting the sides."
// @Router /sides/{id} [delete]
func DeleteSides(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sides id"})
		return
	}
	idUint := uint(idInt)

	err = sidesHandler.DeleteSides(idUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting sides"})
		return
	}

	c.Status(http.StatusNoContent)
}
