package profile_service

import (
	"context"
	"errors"

	"github.com/Android12349/food_recomendation/profile_managment_service/internal/models"
)

func (s *ProfileService) CreateProduct(ctx context.Context, product *models.Product) error {
	// Проверяем существование пользователя
	_, err := s.profileStorage.GetUserByID(ctx, product.UserID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	if err := s.validateProduct(product); err != nil {
		return err
	}
	return s.profileStorage.CreateProduct(ctx, product)
}

func (s *ProfileService) GetProductsByUserID(ctx context.Context, userID int32) ([]*models.Product, error) {
	return s.profileStorage.GetProductsByUserID(ctx, userID)
}

func (s *ProfileService) GetProductByID(ctx context.Context, id int32) (*models.Product, error) {
	return s.profileStorage.GetProductByID(ctx, id)
}

func (s *ProfileService) UpdateProduct(ctx context.Context, product *models.Product) error {
	// Проверяем существование продукта
	_, err := s.profileStorage.GetProductByID(ctx, product.ID)
	if err != nil {
		return errors.New("продукт не найден")
	}

	// Проверяем существование пользователя
	_, err = s.profileStorage.GetUserByID(ctx, product.UserID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	if err := s.validateProduct(product); err != nil {
		return err
	}
	return s.profileStorage.UpdateProduct(ctx, product)
}

func (s *ProfileService) DeleteProduct(ctx context.Context, id int32) error {
	return s.profileStorage.DeleteProduct(ctx, id)
}
