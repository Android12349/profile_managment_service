package bootstrap

import (
	server "github.com/Android12349/food_recomendation/profile_managment_service/internal/api/profile_management_api"
	"github.com/Android12349/food_recomendation/profile_managment_service/internal/services/profile_service"
)

func InitProfileManagementAPI(profileService *profile_service.ProfileService) *server.ProfileManagementAPI {
	return server.NewProfileManagementAPI(profileService)
}
