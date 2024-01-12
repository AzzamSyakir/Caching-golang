package cache

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// RedisClient adalah instance Redis client
var RedisClient *redis.Client
var RedisAddr string // tambahkan variabel untuk menyimpan alamat Redis

// InitRedis inisialisasi koneksi ke Redis
func InitRedis() *redis.Client {
	options := &redis.Options{
		Addr:     "redis:6379",
		Password: "root",
		DB:       1,
	}

	RedisClient = redis.NewClient(options)

	// Uji koneksi ke Redis
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return nil
	}

	// Koneksi ke Redis berhasil, kembalikan RedisClient
	return RedisClient
}

// GetRedisStatus mengembalikan status koneksi Redis
func GetRedisStatus() string {
	// Dapatkan port Redis dari options
	_, portStr, err := net.SplitHostPort(RedisClient.Options().Addr)
	if err != nil {
		return fmt.Sprintf("Error getting Redis port: %v", err)
	}

	// Konversi port string ke integer
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Sprintf("Error converting Redis port to integer: %v", err)
	}

	// Format string dengan alamat dan port Redis
	return fmt.Sprintf("Connected to Redis: %d", port)
}

// CloseRedis menutup koneksi Redis
func CloseRedis() error {
	if err := RedisClient.Close(); err != nil {
		fmt.Println("Error closing Redis client:", err)
		return err
	}
	return nil
}
