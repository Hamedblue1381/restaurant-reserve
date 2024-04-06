package v1

import (
	"net/http"
	"strconv"

	"github.com/Hamedblue1381/restaurant-reserve/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var foodHandler *models.FoodHandler

func InitializedFoodHandler(db *gorm.DB) {
	foodHandler = models.NewFoodHandler(db)
}

// @Summary Get a Single food Dish
// @Description Retrieves details of a single food dish by their unique identifier.
// @Tags food
// @Produce json
// @Param id path int true "food ID" Format(int64)
// @Security Bearer
// @Success 200 {object} models.Food "The details of the food including ID, name, quantity, category, mealtype."
// @Failure 400 {object} ErrorResponse "Invalid food ID format."
// @Failure 404 {object} ErrorResponse "Food not found with the specified ID."
// @Router /food/{id} [get]
func GetFood(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid food id"})
		return
	}

	idUint := uint(idInt)

	user, err := foodHandler.GetFood(idUint)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Food not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Get All Foods
// @Description Retrieves a list of all foods in the system.
// @Tags user
// @Produce json
// @Security Bearer
// @Success 200 {array} models.Food "An array of food objects."
// @Failure 500 {object} ErrorResponse "Internal server error while fetching foods."
// @Router /foods [get]
func GetFoods(c *gin.Context) {
	foods, err := foodHandler.GetFoods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching foods!"})
		return
	}

	c.JSON(http.StatusOK, foods)
}

// @Summary Create a New Food
// @Description Adds a new Food to the system with the provided details.
// @Tags food
// @Accept json
// @Produce json
// @Param food body models.Food true "Food Details"
// @Security Bearer
// @Success 201 {object} models.Food "The created Food's details, including their unique identifier."
// @Failure 400 {object} ErrorResponse "Invalid input format for Food."
// @Failure 500 {object} ErrorResponse "Internal server error while creating the food."
// @Router /food [post]
func CreateFood(c *gin.Context) {
	var food models.Food

	if err := c.ShouldBindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := foodHandler.CreateFood(&food); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating food!"})
		return
	}

	c.JSON(http.StatusCreated, food)
}

// @Summary Update a food
// @Description Updates the details of an existing food identified by their ID.
// @Tags food
// @Accept json
// @Produce json
// @Param id path int true "Food ID" Format(int64)
// @Param food body models.Food true "Updated food Details"
// @Security Bearer
// @Success 200 {object} models.Food "The updated food's details."
// @Failure 400 {object} ErrorResponse "Invalid input format for user details or invalid food ID."
// @Failure 500 {object} ErrorResponse "Internal server error while updating the food."
// @Router /food/{id} [put]
func UpdateFood(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid food id"})
		return
	}
	idUint := uint(idInt)

	var food models.Food

	if err := c.ShouldBindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = foodHandler.UpdateFood(idUint, &food)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating food"})
		return
	}

	c.JSON(http.StatusOK, food)
}

// @Summary Delete a food
// @Description Removes a food dish from the system by their unique identifier.
// @Tags foods
// @Produce json
// @Param id path int true "food ID" Format(int64)
// @Security Bearer
// @Success 204 "Food successfully deleted, no content to return."
// @Failure 400 {object} ErrorResponse "Invalid food ID format."
// @Failure 500 {object} ErrorResponse "Internal server error while deleting the food."
// @Router /food/{id} [delete]
func DeleteFood(c *gin.Context) {
	idString := c.Param("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid food id"})
		return
	}
	idUint := uint(idInt)

	err = foodHandler.DeleteFood(idUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting food"})
		return
	}

	c.Status(http.StatusNoContent)
}
