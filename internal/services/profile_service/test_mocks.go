package profile_service

import (
	"context"
	"errors"

	"github.com/Android12349/food_recomendation/profile_managment_service/internal/models"
)

// mockMenuGenerationProducer - простой mock для тестов
type mockMenuGenerationProducer struct{}

func (m *mockMenuGenerationProducer) PublishMenuGenerationRequest(ctx context.Context, userID int32, bju *models.BJU, budget *int32, preferences string, productNames []string) error {
	// В тестах просто игнорируем публикацию
	return nil
}

// mockMenuGenerationProducerWithError - mock для тестирования ошибок Kafka
type mockMenuGenerationProducerWithError struct{}

func (m *mockMenuGenerationProducerWithError) PublishMenuGenerationRequest(ctx context.Context, userID int32, bju *models.BJU, budget *int32, preferences string, productNames []string) error {
	// Возвращаем ошибку для тестирования обработки ошибок Kafka
	return errors.New("kafka publish error")
}

