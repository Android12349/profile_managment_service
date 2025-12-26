package profile_service

import (
	"context"
	"errors"
	"testing"

	"github.com/Android12349/food_recomendation/profile_managment_service/internal/models"
	"github.com/Android12349/food_recomendation/profile_managment_service/internal/services/profile_service/mocks"
	"github.com/stretchr/testify/suite"
	"gotest.tools/v3/assert"
)

type MealServiceSuite struct {
	suite.Suite
	ctx            context.Context
	profileStorage *mocks.ProfileStorage
	profileService *ProfileService
}

func (s *MealServiceSuite) SetupTest() {
	s.profileStorage = mocks.NewProfileStorage(s.T())
	s.ctx = context.Background()
	mockProducer := &mockMenuGenerationProducer{}
	s.profileService = NewProfileService(s.ctx, s.profileStorage, mockProducer, 3, 50, 6)
}

func (s *MealServiceSuite) TestCreateMealSuccess() {
	meal := testMeal(0, 1, "Курица с рисом", []int32{1, 2})
	user := testUser(1, "testuser")
	product1 := testProduct(1, 1, "Куриная грудка")
	product2 := testProduct(2, 1, "Рис")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, meal.UserID).Return(user, nil)
	s.profileStorage.EXPECT().GetProductByID(s.ctx, int32(1)).Return(product1, nil)
	s.profileStorage.EXPECT().GetProductByID(s.ctx, int32(2)).Return(product2, nil)
	s.profileStorage.EXPECT().CreateMeal(s.ctx, meal).Return(nil)

	got := s.profileService.CreateMeal(s.ctx, meal)
	assert.NilError(s.T(), got)
}

func (s *MealServiceSuite) TestCreateMealUserNotFound() {
	meal := testMeal(0, 1, "Курица с рисом", []int32{1, 2})
	want := errors.New("user not found")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, meal.UserID).Return(nil, want)

	got := s.profileService.CreateMeal(s.ctx, meal)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "пользователь не найден")
}

func (s *MealServiceSuite) TestCreateMealValidationError_EmptyName() {
	meal := testMeal(0, 1, "", []int32{1})
	user := testUser(1, "testuser")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, meal.UserID).Return(user, nil)

	got := s.profileService.CreateMeal(s.ctx, meal)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "название блюда не может быть пустым")
}

func (s *MealServiceSuite) TestCreateMealValidationError_EmptyProductIDs() {
	meal := testMeal(0, 1, "Курица с рисом", []int32{})
	user := testUser(1, "testuser")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, meal.UserID).Return(user, nil)

	got := s.profileService.CreateMeal(s.ctx, meal)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "блюдо должно содержать хотя бы один продукт")
}

func (s *MealServiceSuite) TestCreateMealProductNotFound() {
	meal := testMeal(0, 1, "Курица с рисом", []int32{1, 2})
	user := testUser(1, "testuser")
	product1 := testProduct(1, 1, "Куриная грудка")
	want := errors.New("product not found")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, meal.UserID).Return(user, nil)
	s.profileStorage.EXPECT().GetProductByID(s.ctx, int32(1)).Return(product1, nil)
	s.profileStorage.EXPECT().GetProductByID(s.ctx, int32(2)).Return(nil, want)

	got := s.profileService.CreateMeal(s.ctx, meal)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "продукт с id 2 не найден")
}

func (s *MealServiceSuite) TestGetMealsByUserIDSuccess() {
	userID := int32(1)
	want := []*models.Meal{
		testMeal(1, userID, "Курица с рисом", []int32{1, 2}),
	}

	s.profileStorage.EXPECT().GetMealsByUserID(s.ctx, userID).Return(want, nil)

	got, err := s.profileService.GetMealsByUserID(s.ctx, userID)
	assert.NilError(s.T(), err)
	assert.Equal(s.T(), len(got), 1)
}

func (s *MealServiceSuite) TestGetMealByIDSuccess() {
	mealID := int32(1)
	want := testMeal(mealID, 1, "Курица с рисом", []int32{1, 2})

	s.profileStorage.EXPECT().GetMealByID(s.ctx, mealID).Return(want, nil)

	got, err := s.profileService.GetMealByID(s.ctx, mealID)
	assert.NilError(s.T(), err)
	assert.Equal(s.T(), got.ID, want.ID)
}

func (s *MealServiceSuite) TestUpdateMealSuccess() {
	meal := testMeal(1, 1, "Обновленная курица с рисом", []int32{1, 2})
	existingMeal := testMeal(1, 1, "Курица с рисом", []int32{1})
	user := testUser(1, "testuser")
	product1 := testProduct(1, 1, "Куриная грудка")
	product2 := testProduct(2, 1, "Рис")

	s.profileStorage.EXPECT().GetMealByID(s.ctx, meal.ID).Return(existingMeal, nil)
	s.profileStorage.EXPECT().GetUserByID(s.ctx, meal.UserID).Return(user, nil)
	s.profileStorage.EXPECT().GetProductByID(s.ctx, int32(1)).Return(product1, nil)
	s.profileStorage.EXPECT().GetProductByID(s.ctx, int32(2)).Return(product2, nil)
	s.profileStorage.EXPECT().UpdateMeal(s.ctx, meal).Return(nil)

	got := s.profileService.UpdateMeal(s.ctx, meal)
	assert.NilError(s.T(), got)
}

func (s *MealServiceSuite) TestUpdateMealNotFound() {
	meal := testMeal(1, 1, "Обновленная курица с рисом", []int32{1})
	want := errors.New("meal not found")

	s.profileStorage.EXPECT().GetMealByID(s.ctx, meal.ID).Return(nil, want)

	got := s.profileService.UpdateMeal(s.ctx, meal)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "блюдо не найдено")
}

func (s *MealServiceSuite) TestUpdateMealUserNotFound() {
	meal := testMeal(1, 1, "Обновленная курица с рисом", []int32{1})
	existingMeal := testMeal(1, 1, "Курица с рисом", []int32{1})
	want := errors.New("user not found")

	s.profileStorage.EXPECT().GetMealByID(s.ctx, meal.ID).Return(existingMeal, nil)
	s.profileStorage.EXPECT().GetUserByID(s.ctx, meal.UserID).Return(nil, want)

	got := s.profileService.UpdateMeal(s.ctx, meal)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "пользователь не найден")
}

func (s *MealServiceSuite) TestUpdateMealValidationError_EmptyName() {
	meal := testMeal(1, 1, "", []int32{1})
	existingMeal := testMeal(1, 1, "Курица с рисом", []int32{1})
	user := testUser(1, "testuser")

	s.profileStorage.EXPECT().GetMealByID(s.ctx, meal.ID).Return(existingMeal, nil)
	s.profileStorage.EXPECT().GetUserByID(s.ctx, meal.UserID).Return(user, nil)

	got := s.profileService.UpdateMeal(s.ctx, meal)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "название блюда не может быть пустым")
}

func (s *MealServiceSuite) TestUpdateMealValidationError_EmptyProductIDs() {
	meal := testMeal(1, 1, "Курица с рисом", []int32{})
	existingMeal := testMeal(1, 1, "Курица с рисом", []int32{1})
	user := testUser(1, "testuser")

	s.profileStorage.EXPECT().GetMealByID(s.ctx, meal.ID).Return(existingMeal, nil)
	s.profileStorage.EXPECT().GetUserByID(s.ctx, meal.UserID).Return(user, nil)

	got := s.profileService.UpdateMeal(s.ctx, meal)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "блюдо должно содержать хотя бы один продукт")
}

func (s *MealServiceSuite) TestUpdateMealProductNotFound() {
	meal := testMeal(1, 1, "Обновленная курица с рисом", []int32{1, 2})
	existingMeal := testMeal(1, 1, "Курица с рисом", []int32{1})
	user := testUser(1, "testuser")
	product1 := testProduct(1, 1, "Куриная грудка")
	want := errors.New("product not found")

	s.profileStorage.EXPECT().GetMealByID(s.ctx, meal.ID).Return(existingMeal, nil)
	s.profileStorage.EXPECT().GetUserByID(s.ctx, meal.UserID).Return(user, nil)
	s.profileStorage.EXPECT().GetProductByID(s.ctx, int32(1)).Return(product1, nil)
	s.profileStorage.EXPECT().GetProductByID(s.ctx, int32(2)).Return(nil, want)

	got := s.profileService.UpdateMeal(s.ctx, meal)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "продукт с id 2 не найден")
}

func (s *MealServiceSuite) TestDeleteMealSuccess() {
	mealID := int32(1)

	s.profileStorage.EXPECT().DeleteMeal(s.ctx, mealID).Return(nil)

	got := s.profileService.DeleteMeal(s.ctx, mealID)
	assert.NilError(s.T(), got)
}

func (s *MealServiceSuite) TestDeleteMealError() {
	mealID := int32(1)
	want := errors.New("delete error")

	s.profileStorage.EXPECT().DeleteMeal(s.ctx, mealID).Return(want)

	got := s.profileService.DeleteMeal(s.ctx, mealID)
	assert.ErrorIs(s.T(), got, want)
}

func TestMealServiceSuite(t *testing.T) {
	suite.Run(t, new(MealServiceSuite))
}
