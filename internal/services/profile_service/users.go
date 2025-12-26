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

	// Публикуем событие в Kafka, если есть параметры для генерации меню
	if s.shouldPublishMenuGenerationEvent(user) {
		// Получаем продукты пользователя для передачи в событие
		products, err := s.profileStorage.GetProductsByUserID(ctx, user.ID)
		if err != nil {
			// Логируем ошибку, но продолжаем публикацию без продуктов
		}
		// Инициализируем как пустой slice, а не nil
		productNames := make([]string, 0)
		if products != nil {
			productNames = make([]string, 0, len(products))
			for _, product := range products {
				productNames = append(productNames, product.Name)
			}
		}
		
		if err := s.menuGenerationProducer.PublishMenuGenerationRequest(ctx, user.ID, user.BJU, user.Budget, user.Preferences, productNames); err != nil {
			// Логируем ошибку, но не возвращаем её, т.к. пользователь уже создан
			// В production можно использовать structured logging
		}
	}

	return nil
}

func (s *ProfileService) GetUserByID(ctx context.Context, id int32) (*models.User, error) {
	return s.profileStorage.GetUserByID(ctx, id)
}

func (s *ProfileService) UpdateUser(ctx context.Context, user *models.User) error {
	// Проверяем существование пользователя
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

	// Публикуем событие в Kafka, если есть параметры для генерации меню
	if s.shouldPublishMenuGenerationEvent(user) {
		// Получаем продукты пользователя для передачи в событие
		products, err := s.profileStorage.GetProductsByUserID(ctx, user.ID)
		if err != nil {
			// Логируем ошибку, но продолжаем публикацию без продуктов
		}
		// Инициализируем как пустой slice, а не nil
		productNames := make([]string, 0)
		if products != nil {
			productNames = make([]string, 0, len(products))
			for _, product := range products {
				productNames = append(productNames, product.Name)
			}
		}
		
		if err := s.menuGenerationProducer.PublishMenuGenerationRequest(ctx, user.ID, user.BJU, user.Budget, user.Preferences, productNames); err != nil {
			// Логируем ошибку, но не возвращаем её, т.к. пользователь уже обновлен
			// В production можно использовать structured logging
		}
	}

	return nil
}

// shouldPublishMenuGenerationEvent проверяет, нужно ли публиковать событие для генерации меню
func (s *ProfileService) shouldPublishMenuGenerationEvent(user *models.User) bool {
	// Публикуем событие, если есть хотя бы БЖУ или бюджет
	return user.BJU != nil || user.Budget != nil
}

func (s *ProfileService) DeleteUser(ctx context.Context, id int32) error {
	return s.profileStorage.DeleteUser(ctx, id)
}
