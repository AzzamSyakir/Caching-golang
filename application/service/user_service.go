package service

import (
	"cache-go/application/entities"
	"cache-go/application/repositories"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) *UserService {
	return &UserService{
		UserRepository: userRepository,
	}
}

func (service *UserService) CreateUser(name, email, password string) error {
	// Hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// Generate a UUID for the user
	uuid := uuid.New()
	userID := uuid.String()

	// Current time
	currentTime := time.Now()

	// Save the user to the database using the validated data
	err = service.UserRepository.CreateUser(userID, name, email, string(hashedPassword), currentTime, currentTime)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) FetchUser() ([]entities.User, error) {
	return us.UserRepository.FetchUser()
}
func (service *UserService) UpdateUser(id string, updatedUser entities.User) (entities.User, error) {
	// Business logic/validation goes here

	// Validate if the user exists
	user, err := service.UserRepository.GetUser(id)
	if err != nil {
		return entities.User{}, err
	}
	if user == nil {
		return entities.User{}, errors.New("user not found")
	}

	// Update only non-empty fields
	if updatedUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return entities.User{}, err
		}
		user.Password = string(hashedPassword)
	}

	if updatedUser.Name != "" {
		user.Name = updatedUser.Name
	}

	if updatedUser.Email != "" {
		user.Email = updatedUser.Email
	}

	// Generate a new UUID for the user
	newUUID := uuid.New()
	user.ID = newUUID.String()

	// Set the updated time
	user.UpdatedAt = time.Now()

	// Call repository to update user in the database
	err = service.UserRepository.UpdateUser(*user)
	if err != nil {
		return entities.User{}, err
	}

	// Return only the updated data
	updatedData := entities.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		UpdatedAt: user.UpdatedAt,
	}

	return updatedData, nil
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
func (us *UserService) GetUser(id string) (*entities.User, error) {
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

	// Membuat zona waktu WIB (Asia/Jakarta)
	wib := time.FixedZone("Asia/Jakarta", 7*60*60)

	// Menentukan waktu kedaluwarsa dalam zona waktu WIB
	expirationTimeWIB := time.Now().In(wib).Add(time.Hour)
	claims["exp"] = expirationTimeWIB.Unix() // Token berlaku selama 1 jam dalam zona waktu WIB

	// Menandatangani token dengan secret key
	secretKeyString := os.Getenv("SECRET_KEY")
	secretKey := []byte(secretKeyString)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	// Generate a UUID for the user
	uuid := uuid.New()
	tokensID := uuid.String()

	// Simpan token dan ID pengguna ke dalam tabel tokens
	err = us.UserRepository.SaveToken(tokensID, user.ID, tokenString, expirationTimeWIB.Unix())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (us *UserService) LogoutUser(tokenString string) error {
	type DataUsers struct {
		Username string `json:"username"`
		UserId   string `json:"user_id"`
		jwt.StandardClaims
	}
	// Parse token dan dapatkan claims
	token, err := jwt.ParseWithClaims(tokenString, &DataUsers{}, func(token *jwt.Token) (interface{}, error) {
		// Verifikasi bahwa metode tanda tangan sesuai
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode tanda tangan tidak valid: %v", token.Header["alg"])
		}
		// Kembalikan secret key untuk memverifikasi tanda tangan
		secretKeyString := os.Getenv("SECRET_KEY")
		return []byte(secretKeyString), nil
	})
	if err != nil {
		return fmt.Errorf("gagal mengurai token: %v", err)
	}

	// Mengekstrak klaim dari token
	claims, ok := token.Claims.(*DataUsers)
	if !ok {
		return nil
	}
	userID := claims.UserId
	currentTime := time.Now()

	us.UserRepository.LogoutUser(userID, currentTime)
	return nil
}
