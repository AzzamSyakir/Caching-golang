package routes

import (
	"cache-go/application/controller"
	"cache-go/application/middleware"
	"cache-go/application/repository"
	"cache-go/application/service"
	"cache-go/config"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(db *sql.DB) *mux.Router {
	userRepo := repository.NewInstance(db)
	userService := service.NewInstance(*userRepo)
	userController := controller.NewInstance(*userService)
	router := mux.NewRouter()
	// auth
	router.HandleFunc("/users", userController.CreateUserController).Methods("POST")
	router.HandleFunc("/users/login", userController.LoginUser).Methods("POST")
	// Router untuk rute dengan dua middleware
	protectedRoutes := router.PathPrefix("/").Subrouter()
	protectedRoutes.Use(middleware.AuthMiddleware)
	// users
	protectedRoutes.HandleFunc("/users", userController.FetchUserController).Methods("GET")
	protectedRoutes.HandleFunc("/users/{id}", userController.GetUserController).Methods("GET")
	protectedRoutes.HandleFunc("/users/{id}", userController.UpdateUserController).Methods("PUT")
	protectedRoutes.HandleFunc("/users/{id}", userController.DeleteUser).Methods("DELETE")

	// Add more routes as needed

	return router
}
func RunServer() {
	db := config.InitDB()
	router := SetupRoutes(db)

	// Mulai server HTTP dengan router yang telah dikonfigurasi
	http.Handle("/", router)
	http.ListenAndServe(":9000", nil)
}
