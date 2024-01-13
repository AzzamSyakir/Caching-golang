package repositories

import (
	"cache-go/application/entities"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) CreateUser(id, name, email, hashedPassword string, createdAt, updatedAt time.Time) error {
	createSQL := `
	    INSERT INTO users (id, name, email, password, created_at, updated_at)
	    VALUES (?, ?, ?, ?, CONVERT_TZ(?, '+00:00', '+07:00'), CONVERT_TZ(?, '+00:00', '+07:00'))
	`

	_, err := ur.db.Exec(createSQL, id, name, email, hashedPassword, createdAt.UTC(), updatedAt.UTC())
	return err
}

func (ur *UserRepository) FetchUser() ([]entities.User, error) {
	rows, err := ur.db.Query("SELECT id, name, email, password, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entities.User

	for rows.Next() {
		var user entities.User
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
func (ur *UserRepository) UpdateUser(user entities.User) error {
	query := "UPDATE users SET updated_at = ?, name = ?, email = ?, password = ? WHERE id = ?"
	_, err := ur.db.Exec(query, user.UpdatedAt, user.Name, user.Email, user.Password, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) DeleteUser(id string) error {
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
		return fmt.Errorf("user with ID %s not found", id)
	}

	return nil
}
func (ur *UserRepository) GetUser(id string) (*entities.User, error) {
	user := &entities.User{}

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
func (ur *UserRepository) LoginUser(email string) (*entities.User, error) {
	user := &entities.User{}

	err := ur.db.QueryRow("SELECT id, name, password FROM users WHERE email=?", email).Scan(&user.ID, &user.Name, &user.Password)
	if err != nil {
		return nil, errors.New("password salah atau pengguna tidak ditemukan")
	}

	return user, nil
}
func (ur *UserRepository) SaveToken(TokensID, userID string, token string, expiration int64) error {
	// SQL statement untuk menyimpan token ke dalam tabel tokens
	saveTokenSQL := `
		INSERT INTO tokens (id, user_id, token, created_at, updated_at, expired_at)
		VALUES (?, ?, ?, CONVERT_TZ(?, '+00:00', '+07:00'), CONVERT_TZ(?, '+00:00', '+07:00'), CONVERT_TZ(?, '+00:00', '+07:00'))
	`

	// Mendapatkan waktu sekarang dalam zona waktu UTC
	now := time.Now().UTC()

	// Konversi waktu kedaluwarsa ke zona waktu UTC
	expirationTimeUTC := time.Unix(expiration, 0).UTC()

	// Menjalankan perintah SQL untuk menyimpan token
	_, err := ur.db.Exec(saveTokenSQL, TokensID, userID, token, now, now, expirationTimeUTC)
	if err != nil {
		return fmt.Errorf("gagal menyimpan token: %v", err)
	}

	return nil
}
func (ur *UserRepository) LogoutUser(userId string, updatedAt time.Time) error {
	logoutUserSQL := `
	    UPDATE tokens SET is_revoked=?, updated_at=CONVERT_TZ(?, '+00:00', '+07:00') WHERE user_id=?
	`

	result, err := ur.db.Exec(logoutUserSQL, 1, updatedAt, userId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("userID %s not found", userId)
	}

	return nil
}
