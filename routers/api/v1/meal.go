package v1

import (
	"net/http"
	"strconv"

	"github.com/Hamedblue1381/restaurant-reserve/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var mealtypeHandler *models.MealTypeHandler

func InitializedMealTypeHandler(db *gorm.DB) {
	mealtypeHandler = models.NewMealTypeHandler(db)
}

// @Summary Get a Single mealtype Dish
// @Description Retrieves details of a single mealtype dish by their unique identifier.
// @Tags mealtype
// @Produce json
// @Param id path int true "mealtype ID" Format(int64)
// @Security Bearer
// @Success 200 {object} models.MealType "The details of the mealtype including ID, name, quantity, category, mealtype."
// @Failure 400 {object} ErrorResponse "Invalid mealtype ID format."
// @Failure 404 {object} ErrorResponse "MealType not found with the specified ID."
// @Router /mealtype/{id} [get]
func GetMealType(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mealtype id"})
		return
	}

	idUint := uint(idInt)

	user, err := mealtypeHandler.GetMealType(idUint)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "MealType not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Get All MealTypes
// @Description Retrieves a list of all mealtypes in the system.
// @Tags mealtype
// @Produce json
// @Security Bearer
// @Success 200 {array} models.MealType "An array of mealtype objects."
// @Failure 500 {object} ErrorResponse "Internal server error while fetching mealtypes."
// @Router /mealtypes [get]
func GetMealTypes(c *gin.Context) {
	mealtypes, err := mealtypeHandler.GetMealTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching mealtypes!"})
		return
	}

	c.JSON(http.StatusOK, mealtypes)
}

// @Summary Create a New MealType
// @Description Adds a new MealType to the system with the provided details.
// @Tags mealtype
// @Accept json
// @Produce json
// @Param mealtype body models.MealType true "MealType Details"
// @Security Bearer
// @Success 201 {object} models.MealType "The created MealType's details, including their unique identifier."
// @Failure 400 {object} ErrorResponse "Invalid input format for MealType."
// @Failure 500 {object} ErrorResponse "Internal server error while creating the mealtype."
// @Router /mealtype [post]
func CreateMealType(c *gin.Context) {
	var mealtype models.MealType

	if err := c.ShouldBindJSON(&mealtype); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := mealtypeHandler.CreateMealType(&mealtype); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating mealtype!"})
		return
	}

	c.JSON(http.StatusCreated, mealtype)
}

// @Summary Update a mealtype
// @Description Updates the details of an existing mealtype identified by their ID.
// @Tags mealtype
// @Accept json
// @Produce json
// @Param id path int true "MealType ID" Format(int64)
// @Param mealtype body models.MealType true "Updated mealtype Details"
// @Security Bearer
// @Success 200 {object} models.MealType "The updated mealtype's details."
// @Failure 400 {object} ErrorResponse "Invalid input format for user details or invalid mealtype ID."
// @Failure 500 {object} ErrorResponse "Internal server error while updating the mealtype."
// @Router /mealtype/{id} [put]
func UpdateMealType(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mealtype id"})
		return
	}
	idUint := uint(idInt)

	var mealtype models.MealType

	if err := c.ShouldBindJSON(&mealtype); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = mealtypeHandler.UpdateMealType(idUint, &mealtype)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating mealtype"})
		return
	}

	c.JSON(http.StatusOK, mealtype)
}

// @Summary Delete a mealtype
// @Description Removes a mealtype dish from the system by their unique identifier.
// @Tags mealtype
// @Produce json
// @Param id path int true "mealtype ID" Format(int64)
// @Security Bearer
// @Success 204 "MealType successfully deleted, no content to return."
// @Failure 400 {object} ErrorResponse "Invalid mealtype ID format."
// @Failure 500 {object} ErrorResponse "Internal server error while deleting the mealtype."
// @Router /mealtype/{id} [delete]
func DeleteMealType(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mealtype id"})
		return
	}
	idUint := uint(idInt)

	err = mealtypeHandler.DeleteMealType(idUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting mealtype"})
		return
	}

	c.Status(http.StatusNoContent)
}
