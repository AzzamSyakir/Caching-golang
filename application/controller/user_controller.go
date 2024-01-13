package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	models "cache-go/application/entities"
	"cache-go/application/responses"
	"cache-go/application/service"

	"github.com/gorilla/mux"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(us service.UserService) *UserController {
	return &UserController{userService: us}
}

func (uc *UserController) CreateUserController(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		responses.ErrorResponse(w, "Failed to read user data from the request", http.StatusBadRequest)
		return
	}

	err := uc.userService.CreateUser(user.Name, user.Email, user.Password)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to create user: %v", err)
		responses.ErrorResponse(w, errorMessage, http.StatusInternalServerError)
		return
	}

	currentTime := time.Now()

	// Membuat objek data pengguna untuk dikirim dalam respons
	userData := struct {
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	responses.SuccessResponse(w, "Success", userData, http.StatusCreated)
}
func (uc *UserController) FetchUserController(w http.ResponseWriter, r *http.Request) {
	usersData, err := uc.userService.FetchUser()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get users: %v", err)
		responses.ErrorResponse(w, errorMessage, http.StatusInternalServerError)
		return
	}

	// Membuat objek data pengguna untuk dikirim dalam respons
	var responseData []struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	for _, user := range usersData {
		userData := struct {
			ID        string `json:"id"`
			Username  string `json:"username"`
			Email     string `json:"email"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
		}{
			ID:        user.ID,
			Username:  user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.String(), // Menggunakan String() untuk mendapatkan representasi string dari time.Time
			UpdatedAt: user.UpdatedAt.String(), // Menggunakan String() untuk mendapatkan representasi string dari time.Time
		}

		responseData = append(responseData, userData)
	}

	// Mengembalikan data pengguna sebagai JSON
	w.Header().Set("Content-Type", "application/json")
	responses.SuccessResponse(w, "Success", responseData, http.StatusOK)
}

func (uc *UserController) GetUserController(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan parameter id
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		responses.ErrorResponse(w, "id harus disertakan", http.StatusBadRequest)
		return
	}

	userData, err := uc.userService.GetUser(id)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get user: %v", err)
		responses.ErrorResponse(w, errorMessage, http.StatusInternalServerError)
		return
	}

	// Membuat objek data pengguna untuk dikirim dalam respons
	responseData := struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Username:  userData.Name,
		Email:     userData.Email,
		CreatedAt: userData.CreatedAt.String(), // Menggunakan String() untuk mendapatkan representasi string dari time.Time
		UpdatedAt: userData.UpdatedAt.String(), // Menggunakan String() untuk mendapatkan representasi string dari time.Time
	}

	// Mengembalikan data pengguna sebagai JSON
	w.Header().Set("Content-Type", "application/json")
	responses.SuccessResponse(w, "Success", responseData, http.StatusOK)
}

func (uc *UserController) UpdateUserController(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan parameter id
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		responses.ErrorResponse(w, "id harus disertakan", http.StatusBadRequest)
		return
	}

	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		responses.ErrorResponse(w, "Failed to read user data from the request", http.StatusBadRequest)
		return
	}

	// Update user in the service layer
	_, err := uc.userService.UpdateUser(id, updatedUser)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to update user: %v", err)
		responses.ErrorResponse(w, errorMessage, http.StatusInternalServerError)
		return
	}

	currentTime := time.Now()

	// Membuat objek data pengguna untuk dikirim dalam respons
	updatedData := models.User{
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		UpdatedAt: currentTime,
	}

	// Return response with only updated data
	w.Header().Set("Content-Type", "application/json")
	responses.SuccessResponse(w, "Success", updatedData, http.StatusCreated)
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan parameter id
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		responses.ErrorResponse(w, "id harus disertakan", http.StatusBadRequest)
		return
	}
	// delete user in the service layer
	err = uc.userService.DeleteUser(id)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to Delete user: %v", err)
		responses.ErrorResponse(w, errorMessage, http.StatusInternalServerError)
		return
	}

	responses.OtherResponses(w, "Success delete user", http.StatusCreated)
}
func (uc *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	var requestUser map[string]string

	// Membaca data JSON dari body permintaan
	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		responses.ErrorResponse(w, "Gagal membaca data pengguna dari permintaan", http.StatusBadRequest)
		return
	}

	// Mendapatkan email dan password dari data pengguna
	email, ok := requestUser["email"]
	if !ok {
		responses.ErrorResponse(w, "Email harus diisi", http.StatusBadRequest)
		return
	}

	password, ok := requestUser["password"]
	if !ok {
		responses.ErrorResponse(w, "Password harus diisi", http.StatusBadRequest)
		return
	}

	// Memanggil service untuk melakukan login
	token, err := uc.userService.LoginUser(email, password)
	if err != nil {
		errorMessage := fmt.Sprintf("login failed: %v", err)

		responses.ErrorResponse(w, errorMessage, http.StatusUnauthorized)
		return
	}

	// Mengembalikan token dan pesan sukses
	response := map[string]string{"token": token}
	responses.SuccessResponse(w, "Login berhasil", response, http.StatusOK)
}

func (us *UserController) LogoutUser(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan token dari header "Authorization"
	tokenString := r.Header.Get("Authorization")

	// Membersihkan token dari string "Bearer "
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	// Memanggil fungsi LogoutUser dari service
	err := us.userService.LogoutUser(tokenString)
	if err != nil {
		// Mengatasi kesalahan
		responses.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Mengembalikan token yang baru setelah logout
	responses.OtherResponses(w, "logout berhasil", http.StatusCreated)
}
