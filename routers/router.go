package routers

import (
	"github.com/Hamedblue1381/restaurant-reserve/config"
	docs "github.com/Hamedblue1381/restaurant-reserve/docs"
	"github.com/Hamedblue1381/restaurant-reserve/middleware"
	"github.com/Hamedblue1381/restaurant-reserve/routers/api"
	v1 "github.com/Hamedblue1381/restaurant-reserve/routers/api/v1"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func UseRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(config.CORSMiddleware())
	docs.SwaggerInfo.Title = "Reservation API"
	docs.SwaggerInfo.Description = "This is a server for managing restaurant with  RestAPI build with Go Gin and Gorm"
	docs.SwaggerInfo.BasePath = "/api/v1"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	apiv1 := r.Group("/api/v1")
	auth := apiv1.Group("/auth")
	auth.POST("/signin", api.Login)
	auth.POST("/register", api.Register)
	apiv1.Use(middleware.Auth())
	{
		// for authorized user
		apiv1.GET("/reservations", middleware.Reserve(), v1.GetReservations)
		apiv1.GET("/reservations/:id", middleware.Reserve(), v1.GetReservation)
		apiv1.GET("/food", v1.GetFoods)
		apiv1.GET("/food:id", v1.GetFood)
		apiv1.GET("/sides", v1.GetSides)
		apiv1.GET("/sides/:id", v1.GetSide)
		apiv1.GET("/mealtype", v1.GetMealTypes)
		apiv1.GET("/mealtype/:id", v1.GetMealType)
		apiv1.GET("/users", v1.GetUsers)
		apiv1.GET("/me", v1.GetMe)
		apiv1.GET("/users/:id", v1.GetUser)
		// apiv1.GET("/users/:id/reservations", middleware.Reserve(), v1.GetUserReservations)
		apiv1.POST("/reservations", middleware.Reserve(), v1.CreateReservation)
		apiv1.PUT("/reservations/:id", middleware.Reserve(), v1.UpdateReservation)
		apiv1.DELETE("/reservations/:id", middleware.Reserve(), v1.DeleteReservation)
		apiv1.POST("/users", v1.CreateUser)
		apiv1.PUT("/users/:id", v1.UpdateUser)
		apiv1.DELETE("/users/:id", v1.DeleteUser)
		// for admin
		adminRoutes := apiv1.Group("/")
		adminRoutes.Use(middleware.Admin())
		{
			adminRoutes.POST("/food", v1.CreateFood)
			adminRoutes.PUT("/food/:id", v1.UpdateFood)
			adminRoutes.DELETE("/food/:id", v1.UpdateFood)

			adminRoutes.POST("/sides", v1.CreateSides)
			adminRoutes.PUT("/sides/:id", v1.UpdateSides)
			adminRoutes.DELETE("/sides/:id", v1.DeleteSides)
		}
	}
	return r
}
