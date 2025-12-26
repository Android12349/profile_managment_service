package profile_management_api

import (
	"context"

	"github.com/Android12349/food_recomendation/profile_managment_service/internal/models"
	"github.com/Android12349/food_recomendation/profile_managment_service/internal/pb/profile_management_api"
	"github.com/Android12349/food_recomendation/profile_managment_service/internal/services/profile_service"
)

type profileService interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id int32) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id int32) error
	CreateProduct(ctx context.Context, product *models.Product) error
	GetProductsByUserID(ctx context.Context, userID int32) ([]*models.Product, error)
	GetProductByID(ctx context.Context, id int32) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, id int32) error
	CreateMeal(ctx context.Context, meal *models.Meal) error
	GetMealsByUserID(ctx context.Context, userID int32) ([]*models.Meal, error)
	GetMealByID(ctx context.Context, id int32) (*models.Meal, error)
	UpdateMeal(ctx context.Context, meal *models.Meal) error
	DeleteMeal(ctx context.Context, id int32) error
}

// ProfileManagementAPI реализует grpc ProfileManagementServiceServer
type ProfileManagementAPI struct {
	profile_management_api.UnimplementedProfileManagementServiceServer
	profileService profileService
}

func NewProfileManagementAPI(profileService *profile_service.ProfileService) *ProfileManagementAPI {
	return &ProfileManagementAPI{
		profileService: profileService,
	}
}
