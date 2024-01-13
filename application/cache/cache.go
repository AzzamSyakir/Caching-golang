package cache

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client // RedisClient adalah instance Redis client
var RedisAddr string          // tambahkan variabel untuk menyimpan alamat Redis
var RedisKey string           // RedisKey adalah kunci untuk menyimpan data user di cache

// InitRedis inisialisasi koneksi ke Redis
func InitRedis() *redis.Client {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		os.Exit(1)
	}

	RedisAddr := os.Getenv("REDIS_ADDR")
	RedisPassword := os.Getenv("REDIS_PW")
	RedisDBStr := os.Getenv("REDIS_DB")

	RedisDB, err := strconv.Atoi(RedisDBStr)
	if err != nil {
		fmt.Println("Error converting Redis DB to integer:", err)
		os.Exit(1)
	}

	options := &redis.Options{
		Addr:     RedisAddr,
		Password: RedisPassword,
		DB:       RedisDB,
	}
	RedisClient = redis.NewClient(options)

	// Uji koneksi ke Redis dan otentikasi
	ctx := context.Background()
	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		// Hentikan program atau lakukan penanganan kesalahan sesuai kebutuhan Anda
		return nil
	}

	// Koneksi ke Redis berhasil, print pesan "Connected to Redis"
	fmt.Printf("Connected to Redis: %s\n", RedisAddr)

	// Kembalikan RedisClient
	return RedisClient
}

// CloseRedis menutup koneksi Redis
func CloseRedis() error {
	if err := RedisClient.Close(); err != nil {
		fmt.Println("Error closing Redis client:", err)
		return err
	}
	return nil
}

// GetCached mengambil data dari cache
func GetCached() ([]byte, error) {
	// Ambil alamat Redis dari .env
	RedisAddr := os.Getenv("REDIS_ADDR")
	RedisKey = fmt.Sprintf("%s:cached_data", RedisAddr)
	ctx := context.Background()

	// Cek apakah data ada di cache
	cachedData, err := RedisClient.Get(ctx, RedisKey).Result()
	if err != nil {
		if err == redis.Nil {
			// Key tidak ditemukan di cache
			return nil, nil
		}
		return nil, fmt.Errorf("error getting cached data: %v", err)
	}

	return []byte(cachedData), nil
}

// SetCached menyimpan data ke cache
func SetCached(data []byte) error {
	ctx := context.Background()

	// Simpan data ke cache dengan waktu kadaluarsa
	err := RedisClient.Set(ctx, RedisKey, data, 0).Err()
	if err != nil {
		return fmt.Errorf("error setting data to cache: %v", err)
	}

	return nil
}
