package routes

import (
	"cache-go/application/controller"
	"cache-go/application/middleware"
	"cache-go/application/repositories"
	"cache-go/application/service"
	"cache-go/config"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(db *sql.DB) *mux.Router {
	// Initialize repositories
	userRepository := repositories.NewUserRepository(db)

	// Initialize services
	userService := service.NewUserService(*userRepository)

	// Initialize controllers
	userController := controller.NewUserController(*userService)

	// Create a new router
	router := mux.NewRouter()

	// Protected routes
	protectedRoutes := router.PathPrefix("/").Subrouter()
	protectedRoutes.Use(middleware.AuthMiddleware)

	// Authentication routes
	router.HandleFunc("/users", userController.CreateUserController).Methods("POST")
	router.HandleFunc("/users/login", userController.LoginUser).Methods("POST")
	router.HandleFunc("/users/logout", userController.LogoutUser).Methods("POST")

	// User routes
	protectedRoutes.HandleFunc("/users", userController.FetchUserController).Methods("GET")
	protectedRoutes.HandleFunc("/users/{id}", userController.GetUserController).Methods("GET")
	protectedRoutes.HandleFunc("/users/{id}", userController.UpdateUserController).Methods("PUT")
	protectedRoutes.HandleFunc("/users/{id}", userController.DeleteUser).Methods("DELETE")

	return router
}

func RunServer() {
	db := config.InitDB()
	router := SetupRoutes(db)

	// Mulai server HTTP dengan router yang telah dikonfigurasi
	http.Handle("/", router)
	http.ListenAndServe(":9000", nil)
}
