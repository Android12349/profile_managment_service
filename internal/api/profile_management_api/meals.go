package profile_management_api

import (
	"context"
	"log"

	"github.com/Android12349/food_recomendation/profile_managment_service/internal/models"
	proto_models "github.com/Android12349/food_recomendation/profile_managment_service/internal/pb/models"
	"github.com/Android12349/food_recomendation/profile_managment_service/internal/pb/profile_management_api"
	"github.com/samber/lo"
)

func (s *ProfileManagementAPI) CreateMeal(ctx context.Context, req *profile_management_api.CreateMealRequest) (*profile_management_api.CreateMealResponse, error) {
	log.Printf("Received CreateMeal request for user_id: %d", req.Meal.UserId)

	meal := mapMealCreateModelToModel(req.Meal)

	err := s.profileService.CreateMeal(ctx, meal)
	if err != nil {
		return &profile_management_api.CreateMealResponse{}, err
	}

	return &profile_management_api.CreateMealResponse{
		Meal: mapMealModelToProto(meal),
	}, nil
}

func (s *ProfileManagementAPI) GetMeals(ctx context.Context, req *profile_management_api.GetMealsRequest) (*profile_management_api.GetMealsResponse, error) {
	log.Printf("Received GetMeals request for user_id: %d", req.UserId)

	meals, err := s.profileService.GetMealsByUserID(ctx, req.UserId)
	if err != nil {
		return &profile_management_api.GetMealsResponse{}, err
	}

	return &profile_management_api.GetMealsResponse{
		Meals: lo.Map(meals, func(m *models.Meal, _ int) *proto_models.MealModel {
			return mapMealModelToProto(m)
		}),
	}, nil
}

func (s *ProfileManagementAPI) UpdateMeal(ctx context.Context, req *profile_management_api.UpdateMealRequest) (*profile_management_api.UpdateMealResponse, error) {
	log.Printf("Received UpdateMeal request for ID: %d", req.Id)

	// Get existing meal to preserve user_id
	existingMeal, err := s.profileService.GetMealByID(ctx, req.Id)
	if err != nil {
		return &profile_management_api.UpdateMealResponse{}, err
	}

	meal := mapMealUpdateModelToModel(req.Meal, req.Id, existingMeal.UserID)

	err = s.profileService.UpdateMeal(ctx, meal)
	if err != nil {
		return &profile_management_api.UpdateMealResponse{}, err
	}

	// Get updated meal
	updatedMeal, err := s.profileService.GetMealByID(ctx, req.Id)
	if err != nil {
		return &profile_management_api.UpdateMealResponse{}, err
	}

	return &profile_management_api.UpdateMealResponse{
		Meal: mapMealModelToProto(updatedMeal),
	}, nil
}

func (s *ProfileManagementAPI) DeleteMeal(ctx context.Context, req *profile_management_api.DeleteMealRequest) (*profile_management_api.DeleteMealResponse, error) {
	log.Printf("Received DeleteMeal request for ID: %d", req.Id)

	err := s.profileService.DeleteMeal(ctx, req.Id)
	if err != nil {
		return &profile_management_api.DeleteMealResponse{}, err
	}

	return &profile_management_api.DeleteMealResponse{}, nil
}

func mapMealCreateModelToModel(protoMeal *proto_models.MealCreateModel) *models.Meal {
	return &models.Meal{
		UserID:     protoMeal.UserId,
		Name:       protoMeal.Name,
		ProductIDs: protoMeal.ProductIds,
	}
}

func mapMealUpdateModelToModel(protoMeal *proto_models.MealUpdateModel, id int32, userID int32) *models.Meal {
	return &models.Meal{
		ID:         id,
		UserID:     userID,
		Name:       protoMeal.Name,
		ProductIDs: protoMeal.ProductIds,
	}
}

func mapMealModelToProto(meal *models.Meal) *proto_models.MealModel {
	return &proto_models.MealModel{
		Id:         meal.ID,
		UserId:     meal.UserID,
		Name:       meal.Name,
		ProductIds: meal.ProductIDs,
		CreatedAt:  meal.CreatedAt,
	}
}
