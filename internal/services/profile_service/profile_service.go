package profile_service

import (
	"context"

	"github.com/Android12349/food_recomendation/profile_managment_service/internal/models"
)

type MenuGenerationProducer interface {
	PublishMenuGenerationRequest(ctx context.Context, userID int32, bju *models.BJU, budget *int32, preferences string, productNames []string) error
}

type ProfileStorage interface {
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

type ProfileService struct {
	profileStorage         ProfileStorage
	menuGenerationProducer MenuGenerationProducer
	minUsernameLen         int
	maxUsernameLen         int
	minPasswordLen         int
}

func NewProfileService(ctx context.Context, profileStorage ProfileStorage, menuGenerationProducer MenuGenerationProducer, minUsernameLen, maxUsernameLen, minPasswordLen int) *ProfileService {
	return &ProfileService{
		profileStorage:         profileStorage,
		menuGenerationProducer: menuGenerationProducer,
		minUsernameLen:         minUsernameLen,
		maxUsernameLen:         maxUsernameLen,
		minPasswordLen:         minPasswordLen,
	}
}
