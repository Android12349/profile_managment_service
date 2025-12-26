package bootstrap

import (
	"fmt"
	"log"

	"github.com/Android12349/food_recomendation/profile_managment_service/config"
	"github.com/Android12349/food_recomendation/profile_managment_service/internal/storage/profile_management_storage"
)

func InitPGStorage(cfg *config.Config) *profile_management_storage.ProfileManagementStorage {
	connectionStrings := make([]string, 0, len(cfg.Database.Shards))
	for _, shard := range cfg.Database.Shards {
		connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			shard.Username, shard.Password, shard.Host, shard.Port, shard.DBName, shard.SSLMode)
		connectionStrings = append(connectionStrings, connectionString)
	}

	bucketCount := cfg.Database.BucketCount
	if bucketCount == 0 {
		bucketCount = len(cfg.Database.Shards) * 8
		log.Printf("bucket_count не указан, используем значение по умолчанию: %d", bucketCount)
	}

	storage, err := profile_management_storage.NewProfileManagementStorage(connectionStrings, bucketCount)
	if err != nil {
		log.Panicf("ошибка инициализации БД, %v", err)
		panic(err)
	}
	return storage
}
