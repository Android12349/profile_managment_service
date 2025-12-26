package profile_service

import (
	"context"
	"errors"

	"github.com/Android12349/food_recomendation/profile_managment_service/internal/models"
)

func (s *ProfileService) CreateUser(ctx context.Context, user *models.User) error {
	if err := s.validateUser(user); err != nil {
		return err
	}

	err := s.profileStorage.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	if s.shouldPublishMenuGenerationEvent(user) {
		products, err := s.profileStorage.GetProductsByUserID(ctx, user.ID)
		if err != nil {
			// Логируем ошибку, но продолжаем публикацию без продуктов
		}
		productNames := make([]string, 0)
		if products != nil {
			productNames = make([]string, 0, len(products))
			for _, product := range products {
				productNames = append(productNames, product.Name)
			}
		}

		if err := s.menuGenerationProducer.PublishMenuGenerationRequest(ctx, user.ID, user.BJU, user.Budget, user.Preferences, productNames); err != nil {
			// Логируем ошибку, но не возвращаем её, т.к. пользователь уже создан
		}
	}

	return nil
}

func (s *ProfileService) GetUserByID(ctx context.Context, id int32) (*models.User, error) {
	return s.profileStorage.GetUserByID(ctx, id)
}

func (s *ProfileService) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := s.profileStorage.GetUserByID(ctx, user.ID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	if err := s.validateUser(user); err != nil {
		return err
	}

	err = s.profileStorage.UpdateUser(ctx, user)
	if err != nil {
		return err
	}

	if s.shouldPublishMenuGenerationEvent(user) {
		products, err := s.profileStorage.GetProductsByUserID(ctx, user.ID)
		if err != nil {
			// Логируем ошибку, но продолжаем публикацию без продуктов
		}
		productNames := make([]string, 0)
		if products != nil {
			productNames = make([]string, 0, len(products))
			for _, product := range products {
				productNames = append(productNames, product.Name)
			}
		}

		if err := s.menuGenerationProducer.PublishMenuGenerationRequest(ctx, user.ID, user.BJU, user.Budget, user.Preferences, productNames); err != nil {
			// Логируем ошибку, но не возвращаем её, т.к. пользователь уже обновлен
		}
	}

	return nil
}

func (s *ProfileService) shouldPublishMenuGenerationEvent(user *models.User) bool {
	return user.BJU != nil || user.Budget != nil || user.Preferences != ""
}

func (s *ProfileService) DeleteUser(ctx context.Context, id int32) error {
	return s.profileStorage.DeleteUser(ctx, id)
}
