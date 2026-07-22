package main

import (
	"log"
	"meal-planner-service/pkg/controllers"
	"meal-planner-service/pkg/service"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// didTestSucceed := db.NewDirectDBTest()
	// if didTestSucceed {
	// 	log.Println("Database test succeeded")
	// } else {
	// 	log.Fatal("Database test failed")
	// }

	godotenv.Load()
	jwtSecret := os.Getenv("SUPABASE_JWT_SECRET")
	log.Println("jwtSecret: ", jwtSecret)
	authService, err := service.NewAuthService(jwtSecret)
	if err != nil {
		log.Fatalf("failed to initialize auth service: %v", err)
	}
	handler := controllers.NewHandler(authService)

	http.HandleFunc("/profile", service.AuthMiddleware(authService, handler.Profile))

	// start server
	log.Println("Server running at :8080")
	http.ListenAndServe(":8080", nil)
}
