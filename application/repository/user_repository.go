package repository

import (
	"cache-go/application/models"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewInstance(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) CreateUser(name, email, hashedPassword string, createdAt, updatedAt time.Time) (int64, error) {
	result, err := ur.db.Exec("INSERT INTO users (name, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		name, email, hashedPassword, createdAt, updatedAt)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}
func (ur *UserRepository) FetchUser() ([]models.User, error) {
	rows, err := ur.db.Query("SELECT id, name, email, password, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
func (ur *UserRepository) UpdateUser(id int, name, email, password string, updatedAt time.Time) error {
	result, err := ur.db.Exec("UPDATE users SET name=?, email=?, password=?, updated_at=? WHERE id=?",
		name, email, password, updatedAt, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// Tidak ada baris yang terpengaruh, user dengan ID tersebut tidak ditemukan
		return fmt.Errorf("user with ID %d not found", id)
	}

	return nil
}
func (ur *UserRepository) DeleteUser(id int) error {
	result, err := ur.db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// Tidak ada baris yang terpengaruh, user dengan ID tersebut tidak ditemukan
		return fmt.Errorf("user with ID %d not found", id)
	}

	return nil
}
func (ur *UserRepository) GetUser(id int) (*models.User, error) {
	user := &models.User{}

	err := ur.db.QueryRow("SELECT id, name, email, password, created_at, updated_at FROM users WHERE id=?", id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (ur *UserRepository) LoginUser(email string) (*models.User, error) {
	user := &models.User{}

	err := ur.db.QueryRow("SELECT id, name, password FROM users WHERE email=?", email).Scan(&user.ID, &user.Name, &user.Password)
	if err != nil {
		return nil, errors.New("password salah atau pengguna tidak ditemukan")
	}

	return user, nil
}
