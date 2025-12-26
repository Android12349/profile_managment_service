package profile_service

import (
	"context"
	"errors"
	"testing"

	"github.com/Android12349/food_recomendation/profile_managment_service/internal/models"
	"github.com/Android12349/food_recomendation/profile_managment_service/internal/services/profile_service/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gotest.tools/v3/assert"
)

type UserServiceSuite struct {
	suite.Suite
	ctx            context.Context
	profileStorage *mocks.ProfileStorage
	profileService *ProfileService
}

func (s *UserServiceSuite) SetupTest() {
	s.profileStorage = mocks.NewProfileStorage(s.T())
	s.ctx = context.Background()
	mockProducer := &mockMenuGenerationProducer{}
	s.profileService = NewProfileService(s.ctx, s.profileStorage, mockProducer, 3, 50, 6)
}

func (s *UserServiceSuite) TestCreateUserSuccess() {
	user := testUserWithParams(0, "testuser", int32Ptr(180), int32Ptr(75), int32Ptr(3000), testBJU(100, 70, 250))

	s.profileStorage.EXPECT().CreateUser(s.ctx, user).Return(nil)
	s.profileStorage.EXPECT().GetProductsByUserID(s.ctx, mock.Anything).Return([]*models.Product{}, nil)

	got := s.profileService.CreateUser(s.ctx, user)
	assert.NilError(s.T(), got)
}

func (s *UserServiceSuite) TestCreateUserValidationError_ShortUsername() {
	user := &models.User{
		Username: "ab",
	}

	got := s.profileService.CreateUser(s.ctx, user)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "username должен быть от 3 до 50 символов")
}

func (s *UserServiceSuite) TestCreateUserValidationError_LongUsername() {
	longUsername := make([]byte, 51)
	for i := range longUsername {
		longUsername[i] = 'a'
	}
	user := &models.User{
		Username: string(longUsername),
	}

	got := s.profileService.CreateUser(s.ctx, user)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "username должен быть от 3 до 50 символов")
}

func (s *UserServiceSuite) TestCreateUserValidationError_NegativeHeight() {
	user := testUserWithParams(0, "testuser", int32Ptr(-10), nil, nil, nil)

	got := s.profileService.CreateUser(s.ctx, user)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "рост должен быть положительным числом")
}

func (s *UserServiceSuite) TestCreateUserValidationError_ZeroHeight() {
	user := testUserWithParams(0, "testuser", int32Ptr(0), nil, nil, nil)

	got := s.profileService.CreateUser(s.ctx, user)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "рост должен быть положительным числом")
}

func (s *UserServiceSuite) TestCreateUserValidationError_NegativeWeight() {
	user := testUserWithParams(0, "testuser", nil, int32Ptr(-10), nil, nil)

	got := s.profileService.CreateUser(s.ctx, user)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "вес должен быть положительным числом")
}

func (s *UserServiceSuite) TestCreateUserValidationError_ZeroWeight() {
	user := testUserWithParams(0, "testuser", nil, int32Ptr(0), nil, nil)

	got := s.profileService.CreateUser(s.ctx, user)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "вес должен быть положительным числом")
}

func (s *UserServiceSuite) TestCreateUserValidationError_NegativeBJU_Protein() {
	user := testUserWithParams(0, "testuser", nil, nil, nil, testBJU(-10, 70, 250))

	got := s.profileService.CreateUser(s.ctx, user)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "белок не может быть отрицательным")
}

func (s *UserServiceSuite) TestCreateUserValidationError_NegativeBJU_Fat() {
	user := testUserWithParams(0, "testuser", nil, nil, nil, testBJU(100, -10, 250))

	got := s.profileService.CreateUser(s.ctx, user)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "жиры не могут быть отрицательными")
}

func (s *UserServiceSuite) TestCreateUserValidationError_NegativeBJU_Carbs() {
	user := testUserWithParams(0, "testuser", nil, nil, nil, testBJU(100, 70, -10))

	got := s.profileService.CreateUser(s.ctx, user)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "углеводы не могут быть отрицательными")
}

func (s *UserServiceSuite) TestCreateUserValidationError_NegativeBudget() {
	user := testUserWithParams(0, "testuser", nil, nil, int32Ptr(-100), nil)

	got := s.profileService.CreateUser(s.ctx, user)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "бюджет должен быть положительным числом")
}

func (s *UserServiceSuite) TestCreateUserValidationError_ZeroBudget() {
	user := testUserWithParams(0, "testuser", nil, nil, int32Ptr(0), nil)

	got := s.profileService.CreateUser(s.ctx, user)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "бюджет должен быть положительным числом")
}

func (s *UserServiceSuite) TestCreateUserStorageError() {
	user := testUser(0, "testuser")
	want := errors.New("storage error")

	s.profileStorage.EXPECT().CreateUser(s.ctx, user).Return(want)

	got := s.profileService.CreateUser(s.ctx, user)
	assert.ErrorIs(s.T(), got, want)
}

func (s *UserServiceSuite) TestCreateUserWithoutBJUAndBudget() {
	user := testUser(0, "testuser")

	s.profileStorage.EXPECT().CreateUser(s.ctx, user).Return(nil)

	got := s.profileService.CreateUser(s.ctx, user)
	assert.NilError(s.T(), got)
}

func (s *UserServiceSuite) TestCreateUserKafkaPublishError() {
	user := testUserWithParams(1, "testuser", nil, nil, int32Ptr(3000), testBJU(100, 70, 250))

	s.profileStorage.EXPECT().CreateUser(s.ctx, user).Return(nil)
	s.profileStorage.EXPECT().GetProductsByUserID(s.ctx, user.ID).Return([]*models.Product{}, nil)

	mockProducer := &mockMenuGenerationProducerWithError{}
	s.profileService = NewProfileService(s.ctx, s.profileStorage, mockProducer, 3, 50, 6)

	got := s.profileService.CreateUser(s.ctx, user)
	assert.NilError(s.T(), got)
}

func (s *UserServiceSuite) TestGetUserByIDSuccess() {
	userID := int32(1)
	want := testUser(userID, "testuser")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, userID).Return(want, nil)

	got, err := s.profileService.GetUserByID(s.ctx, userID)
	assert.NilError(s.T(), err)
	assert.Equal(s.T(), got.ID, want.ID)
	assert.Equal(s.T(), got.Username, want.Username)
}

func (s *UserServiceSuite) TestGetUserByIDError() {
	userID := int32(1)
	want := errors.New("user not found")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, userID).Return(nil, want)

	got, err := s.profileService.GetUserByID(s.ctx, userID)
	assert.ErrorIs(s.T(), err, want)
	assert.Check(s.T(), got == nil)
}

func (s *UserServiceSuite) TestUpdateUserSuccess() {
	user := testUser(1, "updateduser")
	user.Height = int32Ptr(185)
	existingUser := testUser(1, "olduser")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, user.ID).Return(existingUser, nil)
	s.profileStorage.EXPECT().UpdateUser(s.ctx, user).Return(nil)

	got := s.profileService.UpdateUser(s.ctx, user)
	assert.NilError(s.T(), got)
}

func (s *UserServiceSuite) TestUpdateUserNotFound() {
	user := testUser(1, "updateduser")
	want := errors.New("user not found")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, user.ID).Return(nil, want)

	got := s.profileService.UpdateUser(s.ctx, user)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "пользователь не найден")
}

func (s *UserServiceSuite) TestUpdateUserValidationError_ShortUsername() {
	user := testUser(1, "ab")
	existingUser := testUser(1, "olduser")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, user.ID).Return(existingUser, nil)

	got := s.profileService.UpdateUser(s.ctx, user)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "username должен быть от 3 до 50 символов")
}

func (s *UserServiceSuite) TestUpdateUserValidationError_NegativeHeight() {
	user := testUserWithParams(1, "testuser", int32Ptr(-10), nil, nil, nil)
	existingUser := testUser(1, "olduser")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, user.ID).Return(existingUser, nil)

	got := s.profileService.UpdateUser(s.ctx, user)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "рост должен быть положительным числом")
}

func (s *UserServiceSuite) TestUpdateUserValidationError_NegativeBJU_Fat() {
	user := testUserWithParams(1, "testuser", nil, nil, nil, testBJU(100, -10, 250))
	existingUser := testUser(1, "olduser")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, user.ID).Return(existingUser, nil)

	got := s.profileService.UpdateUser(s.ctx, user)
	assert.Check(s.T(), got != nil)
	assert.ErrorContains(s.T(), got, "жиры не могут быть отрицательными")
}

func (s *UserServiceSuite) TestDeleteUserSuccess() {
	userID := int32(1)

	s.profileStorage.EXPECT().DeleteUser(s.ctx, userID).Return(nil)

	got := s.profileService.DeleteUser(s.ctx, userID)
	assert.NilError(s.T(), got)
}

func (s *UserServiceSuite) TestDeleteUserError() {
	userID := int32(1)
	want := errors.New("delete error")

	s.profileStorage.EXPECT().DeleteUser(s.ctx, userID).Return(want)

	got := s.profileService.DeleteUser(s.ctx, userID)
	assert.ErrorIs(s.T(), got, want)
}

func (s *UserServiceSuite) TestUpdateUserWithoutBJUAndBudget() {
	user := testUser(1, "updateduser")
	existingUser := testUser(1, "olduser")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, user.ID).Return(existingUser, nil)
	s.profileStorage.EXPECT().UpdateUser(s.ctx, user).Return(nil)

	got := s.profileService.UpdateUser(s.ctx, user)
	assert.NilError(s.T(), got)
}

func (s *UserServiceSuite) TestUpdateUserKafkaPublishError() {
	user := testUserWithParams(1, "testuser", nil, nil, int32Ptr(3000), testBJU(100, 70, 250))
	existingUser := testUser(1, "olduser")

	s.profileStorage.EXPECT().GetUserByID(s.ctx, user.ID).Return(existingUser, nil)
	s.profileStorage.EXPECT().UpdateUser(s.ctx, user).Return(nil)
	s.profileStorage.EXPECT().GetProductsByUserID(s.ctx, user.ID).Return([]*models.Product{}, nil)

	mockProducer := &mockMenuGenerationProducerWithError{}
	s.profileService = NewProfileService(s.ctx, s.profileStorage, mockProducer, 3, 50, 6)

	got := s.profileService.UpdateUser(s.ctx, user)
	assert.NilError(s.T(), got)
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}
