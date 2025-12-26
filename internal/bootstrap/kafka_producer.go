package bootstrap

import (
	"fmt"

	"github.com/Android12349/food_recomendation/profile_managment_service/config"
	"github.com/Android12349/food_recomendation/profile_managment_service/internal/producer/menu_generation_producer"
)

func InitMenuGenerationProducer(cfg *config.Config) *menu_generation_producer.MenuGenerationProducer {
	brokers := []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)}
	return menu_generation_producer.NewMenuGenerationProducer(brokers, cfg.Kafka.MenuGenerationTopicName)
}

