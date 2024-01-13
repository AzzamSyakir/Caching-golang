package cache

import (
	"cache-go/application/entities"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

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

// mengambil data dari cache berdasarkan Redis key
func FetchAllDataFromCache(redisKey string) ([]entities.User, error) {
	ctx := context.Background()

	// Inisialisasi variabel untuk menyimpan hasil scan
	var keys []string
	var cursor uint64

	// Kumpulkan semua keys yang diawali dengan "user:"
	var cachedUsers []entities.User //definisikan var diluar loop
	for {
		var err error
		keys, cursor, err = RedisClient.Scan(ctx, cursor, redisKey+":*", 100).Result()
		if err != nil {
			return nil, fmt.Errorf("error scanning keys: %v", err)
		}

		// Loop melalui keys dan ambil data dari cache

		for _, key := range keys {
			cachedData, err := RedisClient.Get(ctx, key).Bytes()
			if err != nil {
				return nil, fmt.Errorf("error getting cached data for key %s: %v", key, err)
			}
			var cachedUser entities.User
			err = json.Unmarshal(cachedData, &cachedUser)
			if err != nil {
				return nil, fmt.Errorf("error unmarshalling cached data for key %s: %v", key, err)
			}
			cachedUsers = append(cachedUsers, cachedUser)
		}

		// Hentikan loop jika sudah selesai scanning (cursor == 0)
		if cursor == 0 {
			break
		}
	}
	return cachedUsers, nil
}

// menyimpan data ke cache dengan Redis key
func SetCached(redisKey string, data []byte) error {
	ctx := context.Background()

	// Simpan data ke cache dengan waktu kadaluarsa
	err := RedisClient.SetEx(ctx, redisKey, data, time.Second*60).Err()
	if err != nil {
		return fmt.Errorf("error setting data to cache: %v", err)
	}

	return nil
}

// menghapus data dari cache berdasarkan RedisKey dan ID
func DeleteSelectedCached(RedisKey string, id string) error {
	// Pastikan RedisClient sudah diinisialisasi sebelum digunakan
	if RedisClient == nil {
		return fmt.Errorf("not connected to Redis")
	}

	// Hapus data dari cache berdasarkan RedisKey dan ID
	ctx := context.Background()
	cacheKey := fmt.Sprintf("%s:%s", RedisKey, id)
	err := RedisClient.Del(ctx, cacheKey).Err()
	if err != nil {
		return fmt.Errorf("error deleting data from cache: %v", err)
	}

	return nil
}
func DeleteCached(RedisKey string) error {
	// Pastikan RedisClient sudah diinisialisasi sebelum digunakan
	if RedisClient == nil {
		return fmt.Errorf("not connected to Redis")
	}

	// Hapus data dari cache berdasarkan RedisKey dan ID
	ctx := context.Background()
	err := RedisClient.Del(ctx, RedisKey).Err()
	if err != nil {
		return fmt.Errorf("error deleting data from cache: %v", err)
	}

	return nil
}
func GetSelectedCached(redisKey string, id string) ([]byte, error) {
	ctx := context.Background()

	// Cek apakah data ada di cache
	cacheKey := fmt.Sprintf("%s:%s", redisKey, id)
	cachedData, err := RedisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			// Key tidak ditemukan di cache
			return nil, nil
		}
		return nil, fmt.Errorf("error getting cached data: %v", err)
	}

	return []byte(cachedData), nil
}
