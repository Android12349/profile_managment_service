package profile_service

import (
	"context"
	"errors"

	"github.com/Android12349/food_recomendation/profile_managment_service/internal/models"
)

type mockMenuGenerationProducer struct{}

func (m *mockMenuGenerationProducer) PublishMenuGenerationRequest(ctx context.Context, userID int32, bju *models.BJU, budget *int32, preferences string, productNames []string) error {
	return nil
}

type mockMenuGenerationProducerWithError struct{}

func (m *mockMenuGenerationProducerWithError) PublishMenuGenerationRequest(ctx context.Context, userID int32, bju *models.BJU, budget *int32, preferences string, productNames []string) error {
	return errors.New("kafka publish error")
}
