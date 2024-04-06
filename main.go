package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/Hamedblue1381/restaurant-reserve/config"
	"github.com/Hamedblue1381/restaurant-reserve/routers"
	"github.com/Hamedblue1381/restaurant-reserve/routers/api"
	v1 "github.com/Hamedblue1381/restaurant-reserve/routers/api/v1"
)

func main() {
	// Load the env
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error loading .env file")
	}

	// Setup database connection
	db := config.SetupDBConnection()
	if db == nil {
		log.Fatal("Failed to connect to database!")
	}

	v1.InitializeReservationHandler(db)
	v1.InitializedFoodHandler(db)
	v1.InitializedMealTypeHandler(db)
	v1.InitializedSidesHandler(db)
	v1.InitializedUserHandler(db)
	api.InitializedAuthHandler(db)

	// Initialize router
	r := routers.UseRouter()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	fmt.Println("Gracefully shutting down server..., press Ctrl+C again to force shutdown")
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}
