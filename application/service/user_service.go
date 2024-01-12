package service

import (
	"cache-go/application/models"
	"cache-go/application/repository"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepository repository.UserRepository
}

func NewInstance(ur repository.UserRepository) *UserService {
	return &UserService{UserRepository: ur}
}

func (service *UserService) CreateUser(name, email, password string) (int64, error) {
	// Hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// Waktu saat ini
	currentTime := time.Now()

	// Simpan pengguna ke database dengan menggunakan data yang telah Anda validasi
	userID, err := service.UserRepository.CreateUser(name, email, string(hashedPassword), currentTime, currentTime)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (us *UserService) FetchUser() ([]models.User, error) {
	return us.UserRepository.FetchUser()
}
func (service *UserService) UpdateUser(id int, name, email, password string) error {
	// Business logic/validation goes here
	// Hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Waktu saat ini
	currentTime := time.Now()

	// Call repository to update user in the database
	err = service.UserRepository.UpdateUser(id, name, email, string(hashedPassword), currentTime)
	if err != nil {
		return err
	}

	return nil
}
func (service *UserService) DeleteUser(id int) error {
	// Business logic/validation goes here

	// Call repository to delete user in the database
	err := service.UserRepository.DeleteUser(id)
	if err != nil {
		return err
	}

	return nil
}
func (us *UserService) GetUser(id int) (*models.User, error) {
	return us.UserRepository.GetUser(id)
}
func (us *UserService) LoginUser(email string, password string) (string, error) {
	// Mendapatkan informasi pengguna dari repository
	user, err := us.UserRepository.LoginUser(email)
	if err != nil {
		return "", err
	}

	// Membandingkan password yang dimasukkan dengan password yang ada di database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("password salah")
	}

	// Jika login berhasil, buat token JWT
	token := jwt.New(jwt.SigningMethodHS256)

	// Menentukan klaim (claims) token
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["username"] = user.Name
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token berlaku selama 24 jam

	// Menandatangani token dengan secret key
	secretKeyString := os.Getenv("SECRET_KEY")
	secretKey := []byte(secretKeyString)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
