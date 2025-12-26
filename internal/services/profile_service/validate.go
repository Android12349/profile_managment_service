package profile_service

import (
	"errors"
	"fmt"

	"github.com/Android12349/food_recomendation/profile_managment_service/internal/models"
)

func (s *ProfileService) validateUser(user *models.User) error {
	if len(user.Username) < s.minUsernameLen || len(user.Username) > s.maxUsernameLen {
		return fmt.Errorf("username должен быть от %d до %d символов", s.minUsernameLen, s.maxUsernameLen)
	}

	if user.Height != nil && *user.Height <= 0 {
		return errors.New("рост должен быть положительным числом")
	}

	if user.Weight != nil && *user.Weight <= 0 {
		return errors.New("вес должен быть положительным числом")
	}

	if user.BJU != nil {
		if user.BJU.Protein < 0 {
			return errors.New("белок не может быть отрицательным")
		}
		if user.BJU.Fat < 0 {
			return errors.New("жиры не могут быть отрицательными")
		}
		if user.BJU.Carbs < 0 {
			return errors.New("углеводы не могут быть отрицательными")
		}
	}

	if user.Budget != nil && *user.Budget <= 0 {
		return errors.New("бюджет должен быть положительным числом")
	}

	return nil
}

func (s *ProfileService) validateProduct(product *models.Product) error {
	if product.Name == "" {
		return errors.New("название продукта не может быть пустым")
	}

	if product.Calories != nil && *product.Calories < 0 {
		return errors.New("калории не могут быть отрицательными")
	}

	if product.Protein != nil && *product.Protein < 0 {
		return errors.New("белок не может быть отрицательным")
	}

	if product.Fat != nil && *product.Fat < 0 {
		return errors.New("жиры не могут быть отрицательными")
	}

	if product.Carbs != nil && *product.Carbs < 0 {
		return errors.New("углеводы не могут быть отрицательными")
	}

	return nil
}

func (s *ProfileService) validateMeal(meal *models.Meal) error {
	if meal.Name == "" {
		return errors.New("название блюда не может быть пустым")
	}

	if len(meal.ProductIDs) == 0 {
		return errors.New("блюдо должно содержать хотя бы один продукт")
	}

	return nil
}
