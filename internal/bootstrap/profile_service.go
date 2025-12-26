package bootstrap

import (
	"context"

	"github.com/Android12349/food_recomendation/profile_managment_service/config"
	"github.com/Android12349/food_recomendation/profile_managment_service/internal/producer/menu_generation_producer"
	"github.com/Android12349/food_recomendation/profile_managment_service/internal/services/profile_service"
	"github.com/Android12349/food_recomendation/profile_managment_service/internal/storage/profile_management_storage"
)

func InitProfileService(storage *profile_management_storage.ProfileManagementStorage, producer *menu_generation_producer.MenuGenerationProducer, cfg *config.Config) *profile_service.ProfileService {
	return profile_service.NewProfileService(
		context.Background(),
		storage,
		producer,
		cfg.ProfileServiceSettings.MinUsernameLen,
		cfg.ProfileServiceSettings.MaxUsernameLen,
		cfg.ProfileServiceSettings.MinPasswordLen,
	)
}
