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

type ProductServiceSuite struct {
	suite.Suite
	ctx            context.Context
	profileStorage *mocks.ProfileStorage
	profileService *ProfileService
}

func (s *ProductServiceSuite) SetupTest() {
	s.profileStorage = mocks.NewProfileStorage(s.T())
	s.ctx = context.Background()
	mockProducer := &mockMenuGenerationProducer{}
	s.profileService = NewProfileService(s.ctx, s.profileStorage, mockProducer, 3, 50, 6)
}

func (s *ProductServiceSuite) TestCreateProductSuccess() {
	product := testProductWithParams(0, 1, "Куриная грудка", int32Ptr(165), int32Ptr(31), int32Ptr(3), int32Ptr(0))
	user := testUser(1, "testuser")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, product.UserID).Return(user, nil)
	s.profileStorage.EXPECT().CreateProduct(s.ctx, product).Return(nil)

	got := s.profileService.CreateProduct(s.ctx, product)
	assert.NilError(s.T(), got)
}

func (s *ProductServiceSuite) TestCreateProductUserNotFound() {
	product := testProduct(0, 1, "Куриная грудка")
	want := errors.New("user not found")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, product.UserID).Return(nil, want)

	got := s.profileService.CreateProduct(s.ctx, product)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "пользователь не найден")
}

func (s *ProductServiceSuite) TestCreateProductValidationError_EmptyName() {
	product := testProduct(0, 1, "")
	user := testUser(1, "testuser")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, product.UserID).Return(user, nil)

	got := s.profileService.CreateProduct(s.ctx, product)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "название продукта не может быть пустым")
}

func (s *ProductServiceSuite) TestCreateProductValidationError_NegativeCalories() {
	product := testProductWithParams(0, 1, "Куриная грудка", int32Ptr(-10), nil, nil, nil)
	user := testUser(1, "testuser")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, product.UserID).Return(user, nil)

	got := s.profileService.CreateProduct(s.ctx, product)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "калории не могут быть отрицательными")
}

func (s *ProductServiceSuite) TestCreateProductValidationError_NegativeProtein() {
	product := testProductWithParams(0, 1, "Куриная грудка", nil, int32Ptr(-10), nil, nil)
	user := testUser(1, "testuser")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, product.UserID).Return(user, nil)

	got := s.profileService.CreateProduct(s.ctx, product)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "белок не может быть отрицательным")
}

func (s *ProductServiceSuite) TestCreateProductValidationError_NegativeFat() {
	product := testProductWithParams(0, 1, "Куриная грудка", nil, nil, int32Ptr(-10), nil)
	user := testUser(1, "testuser")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, product.UserID).Return(user, nil)

	got := s.profileService.CreateProduct(s.ctx, product)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "жиры не могут быть отрицательными")
}

func (s *ProductServiceSuite) TestCreateProductValidationError_NegativeCarbs() {
	product := testProductWithParams(0, 1, "Куриная грудка", nil, nil, nil, int32Ptr(-10))
	user := testUser(1, "testuser")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, product.UserID).Return(user, nil)

	got := s.profileService.CreateProduct(s.ctx, product)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "углеводы не могут быть отрицательными")
}

func (s *ProductServiceSuite) TestCreateProductStorageError() {
	product := testProduct(0, 1, "Куриная грудка")
	user := testUser(1, "testuser")
	want := errors.New("storage error")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, product.UserID).Return(user, nil)
	s.profileStorage.EXPECT().CreateProduct(s.ctx, product).Return(want)

	got := s.profileService.CreateProduct(s.ctx, product)
	assert.ErrorIs(s.T(), got, want)
}

func (s *ProductServiceSuite) TestGetProductsByUserIDSuccess() {
	userID := int32(1)
	want := []*models.Product{
		testProduct(1, userID, "Куриная грудка"),
		testProduct(2, userID, "Рис"),
	}

	s.profileStorage.EXPECT().GetProductsByUserID(s.ctx, userID).Return(want, nil)

	got, err := s.profileService.GetProductsByUserID(s.ctx, userID)
	assert.NilError(s.T(), err)
	assert.Equal(s.T(), len(got), 2)
}

func (s *ProductServiceSuite) TestGetProductByIDSuccess() {
	productID := int32(1)
	want := testProduct(productID, 1, "Куриная грудка")

	s.profileStorage.EXPECT().GetProductByID(s.ctx, productID).Return(want, nil)

	got, err := s.profileService.GetProductByID(s.ctx, productID)
	assert.NilError(s.T(), err)
	assert.Equal(s.T(), got.ID, want.ID)
}

func (s *ProductServiceSuite) TestUpdateProductSuccess() {
	product := testProduct(1, 1, "Обновленная куриная грудка")
	existingProduct := testProduct(1, 1, "Куриная грудка")
	user := testUser(1, "testuser")

	s.profileStorage.EXPECT().GetProductByID(s.ctx, product.ID).Return(existingProduct, nil)
	s.profileStorage.EXPECT().GetUserByID(s.ctx, product.UserID).Return(user, nil)
	s.profileStorage.EXPECT().UpdateProduct(s.ctx, product).Return(nil)

	got := s.profileService.UpdateProduct(s.ctx, product)
	assert.NilError(s.T(), got)
}

func (s *ProductServiceSuite) TestUpdateProductNotFound() {
	product := testProduct(1, 1, "Обновленная куриная грудка")
	want := errors.New("product not found")

	s.profileStorage.EXPECT().GetProductByID(s.ctx, product.ID).Return(nil, want)

	got := s.profileService.UpdateProduct(s.ctx, product)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "продукт не найден")
}

func (s *ProductServiceSuite) TestUpdateProductUserNotFound() {
	product := testProduct(1, 1, "Обновленная куриная грудка")
	existingProduct := testProduct(1, 1, "Куриная грудка")
	want := errors.New("user not found")

	s.profileStorage.EXPECT().GetProductByID(s.ctx, product.ID).Return(existingProduct, nil)
	s.profileStorage.EXPECT().GetUserByID(s.ctx, product.UserID).Return(nil, want)

	got := s.profileService.UpdateProduct(s.ctx, product)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "пользователь не найден")
}

func (s *ProductServiceSuite) TestUpdateProductValidationError_EmptyName() {
	product := testProduct(1, 1, "")
	existingProduct := testProduct(1, 1, "Куриная грудка")
	user := testUser(1, "testuser")

	s.profileStorage.EXPECT().GetProductByID(s.ctx, product.ID).Return(existingProduct, nil)
	s.profileStorage.EXPECT().GetUserByID(s.ctx, product.UserID).Return(user, nil)

	got := s.profileService.UpdateProduct(s.ctx, product)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "название продукта не может быть пустым")
}

func (s *ProductServiceSuite) TestUpdateProductValidationError_NegativeCalories() {
	product := testProductWithParams(1, 1, "Куриная грудка", int32Ptr(-10), nil, nil, nil)
	existingProduct := testProduct(1, 1, "Куриная грудка")
	user := testUser(1, "testuser")

	s.profileStorage.EXPECT().GetProductByID(s.ctx, product.ID).Return(existingProduct, nil)
	s.profileStorage.EXPECT().GetUserByID(s.ctx, product.UserID).Return(user, nil)

	got := s.profileService.UpdateProduct(s.ctx, product)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "калории не могут быть отрицательными")
}

func (s *ProductServiceSuite) TestDeleteProductSuccess() {
	productID := int32(1)

	s.profileStorage.EXPECT().DeleteProduct(s.ctx, productID).Return(nil)

	got := s.profileService.DeleteProduct(s.ctx, productID)
	assert.NilError(s.T(), got)
}

func (s *ProductServiceSuite) TestDeleteProductError() {
	productID := int32(1)
	want := errors.New("delete error")

	s.profileStorage.EXPECT().DeleteProduct(s.ctx, productID).Return(want)

	got := s.profileService.DeleteProduct(s.ctx, productID)
	assert.ErrorIs(s.T(), got, want)
}

func TestProductServiceSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceSuite))
}
