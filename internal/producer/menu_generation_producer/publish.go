package menu_generation_producer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Android12349/food_recomendation/profile_managment_service/internal/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

func (p *MenuGenerationProducer) PublishMenuGenerationRequest(ctx context.Context, userID int32, bju *models.BJU, budget *int32, preferences string, productNames []string) error {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(p.kafkaBroker...),
		Topic:    p.topicName,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	// Создаем событие
	event := models.MenuGenerationRequestEvent{
		RequestID: uuid.New().String(),
		UserID:    userID,
		Preferences: models.MenuGenerationPrefs{
			BJU:      bju,
			Budget:   budget,
			Products: productNames, // Используем продукты из БД
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	// Дополнительно парсим preferences (JSON строка) для извлечения продуктов, если они там есть
	// Но приоритет у продуктов из БД
	if len(productNames) == 0 && preferences != "" {
		var prefsMap map[string]interface{}
		if err := json.Unmarshal([]byte(preferences), &prefsMap); err == nil {
			if products, ok := prefsMap["products"].([]interface{}); ok {
				productStrings := make([]string, 0, len(products))
				for _, prod := range products {
					if str, ok := prod.(string); ok {
						productStrings = append(productStrings, str)
					}
				}
				event.Preferences.Products = productStrings
			}
		}
	}

	// Сериализуем событие в JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "failed to marshal menu generation event")
	}

	// Публикуем сообщение в Kafka
	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("user_%d", userID)),
		Value: eventJSON,
	}

	err = writer.WriteMessages(ctx, msg)
	if err != nil {
		return errors.Wrap(err, "failed to write message to kafka")
	}

	return nil
}
