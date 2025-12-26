package profile_service

import (
	"github.com/Android12349/food_recomendation/profile_managment_service/internal/models"
)

// int32Ptr создает указатель на int32
func int32Ptr(v int32) *int32 {
	return &v
}

// testUser создает тестового пользователя с заданными параметрами
func testUser(id int32, username string) *models.User {
	return &models.User{
		ID:       id,
		Username: username,
	}
}

// testUserWithParams создает тестового пользователя с дополнительными параметрами
func testUserWithParams(id int32, username string, height, weight, budget *int32, bju *models.BJU) *models.User {
	user := testUser(id, username)
	user.Height = height
	user.Weight = weight
	user.Budget = budget
	user.BJU = bju
	return user
}

// testBJU создает тестовый объект БЖУ
func testBJU(protein, fat, carbs int32) *models.BJU {
	return &models.BJU{
		Protein: protein,
		Fat:     fat,
		Carbs:   carbs,
	}
}

// testProduct создает тестовый продукт с заданными параметрами
func testProduct(id, userID int32, name string) *models.Product {
	return &models.Product{
		ID:     id,
		UserID: userID,
		Name:   name,
	}
}

// testProductWithParams создает тестовый продукт с дополнительными параметрами
func testProductWithParams(id, userID int32, name string, calories, protein, fat, carbs *int32) *models.Product {
	product := testProduct(id, userID, name)
	product.Calories = calories
	product.Protein = protein
	product.Fat = fat
	product.Carbs = carbs
	return product
}

// testMeal создает тестовое блюдо с заданными параметрами
func testMeal(id, userID int32, name string, productIDs []int32) *models.Meal {
	return &models.Meal{
		ID:         id,
		UserID:     userID,
		Name:       name,
		ProductIDs: productIDs,
	}
}
