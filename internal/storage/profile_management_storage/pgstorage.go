package profile_management_storage

import (
	"context"
	"fmt"
	"hash/fnv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type ProfileManagementStorage struct {
	shards        []*pgxpool.Pool
	bucketCount   int
	bucketToShard []int
}

func NewProfileManagementStorage(connStrings []string, bucketCount int) (*ProfileManagementStorage, error) {
	if len(connStrings) == 0 {
		return nil, errors.New("необходимо указать хотя бы один шард")
	}

	if bucketCount < len(connStrings) {
		return nil, errors.New("количество бакетов должно быть >= количества шардов")
	}

	shards := make([]*pgxpool.Pool, 0, len(connStrings))
	for i, connString := range connStrings {
		config, err := pgxpool.ParseConfig(connString)
		if err != nil {
			return nil, errors.Wrapf(err, "ошибка парсинга конфига для шарда %d", i)
		}

		db, err := pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			return nil, errors.Wrapf(err, "ошибка подключения к шарду %d", i)
		}
		shards = append(shards, db)
	}

	bucketToShard := make([]int, bucketCount)
	for i := 0; i < bucketCount; i++ {
		bucketToShard[i] = i % len(shards)
	}

	storage := &ProfileManagementStorage{
		shards:        shards,
		bucketCount:   bucketCount,
		bucketToShard: bucketToShard,
	}

	err := storage.initTables()
	if err != nil {
		return nil, err
	}

	return storage, nil
}

// getBucket вычисляет номер бакета для user_id используя хеш-функцию
func (s *ProfileManagementStorage) getBucket(userID int32) int {
	// Используем FNV-1a хеш для равномерного распределения
	hash := fnv.New32a()
	hash.Write([]byte(fmt.Sprintf("%d", userID)))
	hashValue := hash.Sum32()

	bucket := int(hashValue) % s.bucketCount
	if bucket < 0 {
		bucket = -bucket
	}
	return bucket
}

func (s *ProfileManagementStorage) getShard(userID int32) *pgxpool.Pool {
	bucket := s.getBucket(userID)
	shardIndex := s.bucketToShard[bucket]
	return s.shards[shardIndex]
}

func (s *ProfileManagementStorage) initTables() error {
	usersSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			%s SERIAL PRIMARY KEY,
			%s VARCHAR(50) UNIQUE NOT NULL,
			%s VARCHAR(255) NOT NULL,
			%s INT,
			%s INT,
			%s JSONB,
			%s INT,
			%s JSONB,
			%s TIMESTAMP DEFAULT NOW()
		)`, usersTableName, usersIDColumn, usersUsernameColumn, usersPasswordHashColumn,
		usersHeightColumn, usersWeightColumn, usersBJUColumn, usersBudgetColumn,
		usersPreferencesColumn, usersCreatedAtColumn)

	productsSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			%s SERIAL PRIMARY KEY,
			%s INT NOT NULL,
			%s VARCHAR(100) NOT NULL,
			%s INT,
			%s INT,
			%s INT,
			%s INT,
			%s TIMESTAMP DEFAULT NOW()
		)`, productsTableName, productsIDColumn, productsUserIDColumn,
		productsNameColumn, productsCaloriesColumn, productsProteinColumn,
		productsFatColumn, productsCarbsColumn, productsCreatedAtColumn)

	mealsSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			%s SERIAL PRIMARY KEY,
			%s INT NOT NULL,
			%s VARCHAR(100) NOT NULL,
			%s INT[],
			%s TIMESTAMP DEFAULT NOW()
		)`, mealsTableName, mealsIDColumn, mealsUserIDColumn,
		mealsNameColumn, mealsProductIDsColumn, mealsCreatedAtColumn)

	for i, shard := range s.shards {
		_, err := shard.Exec(context.Background(), usersSQL)
		if err != nil {
			return errors.Wrapf(err, "init users table on shard %d", i)
		}

		_, err = shard.Exec(context.Background(), productsSQL)
		if err != nil {
			return errors.Wrapf(err, "init products table on shard %d", i)
		}

		_, err = shard.Exec(context.Background(), mealsSQL)
		if err != nil {
			return errors.Wrapf(err, "init meals table on shard %d", i)
		}
	}

	return nil
}
