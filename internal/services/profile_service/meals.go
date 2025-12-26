package profile_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Android12349/food_recomendation/profile_managment_service/internal/models"
)

func (s *ProfileService) CreateMeal(ctx context.Context, meal *models.Meal) error {
	// Проверяем существование пользователя
	_, err := s.profileStorage.GetUserByID(ctx, meal.UserID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	if err := s.validateMeal(meal); err != nil {
		return err
	}

	// Проверяем существование всех продуктов
	for _, productID := range meal.ProductIDs {
		_, err := s.profileStorage.GetProductByID(ctx, productID)
		if err != nil {
			return fmt.Errorf("продукт с id %d не найден", productID)
		}
	}

	return s.profileStorage.CreateMeal(ctx, meal)
}

func (s *ProfileService) GetMealsByUserID(ctx context.Context, userID int32) ([]*models.Meal, error) {
	return s.profileStorage.GetMealsByUserID(ctx, userID)
}

func (s *ProfileService) GetMealByID(ctx context.Context, id int32) (*models.Meal, error) {
	return s.profileStorage.GetMealByID(ctx, id)
}

func (s *ProfileService) UpdateMeal(ctx context.Context, meal *models.Meal) error {
	// Проверяем существование блюда
	_, err := s.profileStorage.GetMealByID(ctx, meal.ID)
	if err != nil {
		return errors.New("блюдо не найдено")
	}

	// Проверяем существование пользователя
	_, err = s.profileStorage.GetUserByID(ctx, meal.UserID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	if err := s.validateMeal(meal); err != nil {
		return err
	}

	// Проверяем существование всех продуктов
	for _, productID := range meal.ProductIDs {
		_, err := s.profileStorage.GetProductByID(ctx, productID)
		if err != nil {
			return fmt.Errorf("продукт с id %d не найден", productID)
		}
	}

	return s.profileStorage.UpdateMeal(ctx, meal)
}

func (s *ProfileService) DeleteMeal(ctx context.Context, id int32) error {
	return s.profileStorage.DeleteMeal(ctx, id)
}
