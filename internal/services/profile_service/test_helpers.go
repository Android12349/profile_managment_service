package profile_service

import (
	"github.com/Android12349/food_recomendation/profile_managment_service/internal/models"
)

func int32Ptr(v int32) *int32 {
	return &v
}

func testUser(id int32, username string) *models.User {
	return &models.User{
		ID:       id,
		Username: username,
	}
}

func testUserWithParams(id int32, username string, height, weight, budget *int32, bju *models.BJU) *models.User {
	user := testUser(id, username)
	user.Height = height
	user.Weight = weight
	user.Budget = budget
	user.BJU = bju
	return user
}

func testBJU(protein, fat, carbs int32) *models.BJU {
	return &models.BJU{
		Protein: protein,
		Fat:     fat,
		Carbs:   carbs,
	}
}

func testProduct(id, userID int32, name string) *models.Product {
	return &models.Product{
		ID:     id,
		UserID: userID,
		Name:   name,
	}
}

func testProductWithParams(id, userID int32, name string, calories, protein, fat, carbs *int32) *models.Product {
	product := testProduct(id, userID, name)
	product.Calories = calories
	product.Protein = protein
	product.Fat = fat
	product.Carbs = carbs
	return product
}

func testMeal(id, userID int32, name string, productIDs []int32) *models.Meal {
	return &models.Meal{
		ID:         id,
		UserID:     userID,
		Name:       name,
		ProductIDs: productIDs,
	}
}
