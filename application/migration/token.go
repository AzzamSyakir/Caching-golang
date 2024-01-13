package migration

import (
	"database/sql"
	"log"
)

// TokenMigrate digunakan untuk menjalankan migrasi tabel token.
func TokenMigrate(db *sql.DB) {
	// SQL statement untuk memeriksa apakah tabel token sudah ada
	checkTableSQL := `
        SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'tokens'
    `

	// Menjalankan perintah SQL untuk memeriksa apakah tabel sudah ada
	var tableCount int
	err := db.QueryRow(checkTableSQL).Scan(&tableCount)
	if err != nil {
		// Menangani kesalahan jika terjadi kesalahan saat memeriksa tabel
		log.Fatal(err)
		return
	}

	if tableCount > 0 {
		// Jika tabel sudah ada, tampilkan pesan
		log.Println("Tabel tokens sudah di migrasi")
		return
	}

	// SQL statement untuk membuat tabel token dengan kolom "revoke"
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS tokens (
            id CHAR(36) NOT NULL PRIMARY KEY,
		user_id CHAR(36),
            token VARCHAR(255) NOT NULL,
            created_at TIMESTAMP,
            updated_at TIMESTAMP,
            expired_at TIMESTAMP,
            is_revoked TINYINT(1) DEFAULT 0,
            FOREIGN KEY (user_id) REFERENCES users(id)
        )
    `

	// Menjalankan perintah SQL untuk membuat tabel
	_, err = db.Exec(createTableSQL)
	if err != nil {
		// Menangani kesalahan jika terjadi kesalahan saat migrasi
		log.Fatal(err)
		return
	}

	// Pesan sukses jika migrasi berhasil
	log.Println("Migrasi tabel tokens berhasil")
}
