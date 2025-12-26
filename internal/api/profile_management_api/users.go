package profile_management_api

import (
	"context"
	"log"

	"github.com/Android12349/food_recomendation/profile_managment_service/internal/models"
	proto_models "github.com/Android12349/food_recomendation/profile_managment_service/internal/pb/models"
	"github.com/Android12349/food_recomendation/profile_managment_service/internal/pb/profile_management_api"
	"golang.org/x/crypto/bcrypt"
)

func (s *ProfileManagementAPI) CreateUser(ctx context.Context, req *profile_management_api.CreateUserRequest) (*profile_management_api.CreateUserResponse, error) {
	log.Printf("Received CreateUser request for username: %s", req.User.Username)

	user := mapUserCreateModelToModel(req.User)

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return &profile_management_api.CreateUserResponse{}, err
	}
	user.PasswordHash = string(hashedPassword)

	err = s.profileService.CreateUser(ctx, user)
	if err != nil {
		return &profile_management_api.CreateUserResponse{}, err
	}

	return &profile_management_api.CreateUserResponse{
		User: mapUserModelToProto(user),
	}, nil
}

func (s *ProfileManagementAPI) GetUser(ctx context.Context, req *profile_management_api.GetUserRequest) (*profile_management_api.GetUserResponse, error) {
	log.Printf("Received GetUser request for ID: %d", req.Id)

	user, err := s.profileService.GetUserByID(ctx, req.Id)
	if err != nil {
		return &profile_management_api.GetUserResponse{}, err
	}

	return &profile_management_api.GetUserResponse{
		User: mapUserModelToProto(user),
	}, nil
}

func (s *ProfileManagementAPI) UpdateUser(ctx context.Context, req *profile_management_api.UpdateUserRequest) (*profile_management_api.UpdateUserResponse, error) {
	log.Printf("Received UpdateUser request for ID: %d", req.Id)

	user := mapUserUpdateModelToModel(req.User, req.Id)

	err := s.profileService.UpdateUser(ctx, user)
	if err != nil {
		return &profile_management_api.UpdateUserResponse{}, err
	}

	// Get updated user
	updatedUser, err := s.profileService.GetUserByID(ctx, req.Id)
	if err != nil {
		return &profile_management_api.UpdateUserResponse{}, err
	}

	return &profile_management_api.UpdateUserResponse{
		User: mapUserModelToProto(updatedUser),
	}, nil
}

func (s *ProfileManagementAPI) DeleteUser(ctx context.Context, req *profile_management_api.DeleteUserRequest) (*profile_management_api.DeleteUserResponse, error) {
	log.Printf("Received DeleteUser request for ID: %d", req.Id)

	err := s.profileService.DeleteUser(ctx, req.Id)
	if err != nil {
		return &profile_management_api.DeleteUserResponse{}, err
	}

	return &profile_management_api.DeleteUserResponse{}, nil
}

func mapUserCreateModelToModel(protoUser *proto_models.UserCreateModel) *models.User {
	user := &models.User{
		Username:    protoUser.Username,
		Preferences: protoUser.Preferences,
	}

	if protoUser.Height != 0 {
		height := protoUser.Height
		user.Height = &height
	}
	if protoUser.Weight != 0 {
		weight := protoUser.Weight
		user.Weight = &weight
	}
	if protoUser.Budget != 0 {
		budget := protoUser.Budget
		user.Budget = &budget
	}
	if protoUser.Bju != nil {
		user.BJU = &models.BJU{
			Protein: protoUser.Bju.Protein,
			Fat:     protoUser.Bju.Fat,
			Carbs:   protoUser.Bju.Carbs,
		}
	}

	return user
}

func mapUserUpdateModelToModel(protoUser *proto_models.UserUpdateModel, id int32) *models.User {
	user := &models.User{
		ID:          id,
		Username:    protoUser.Username,
		Preferences: protoUser.Preferences,
	}

	if protoUser.Height != 0 {
		height := protoUser.Height
		user.Height = &height
	}
	if protoUser.Weight != 0 {
		weight := protoUser.Weight
		user.Weight = &weight
	}
	if protoUser.Budget != 0 {
		budget := protoUser.Budget
		user.Budget = &budget
	}
	if protoUser.Bju != nil {
		user.BJU = &models.BJU{
			Protein: protoUser.Bju.Protein,
			Fat:     protoUser.Bju.Fat,
			Carbs:   protoUser.Bju.Carbs,
		}
	}

	return user
}

func mapUserModelToProto(user *models.User) *proto_models.UserModel {
	protoUser := &proto_models.UserModel{
		Id:           user.ID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Preferences:  user.Preferences,
		CreatedAt:    user.CreatedAt,
	}

	if user.Height != nil {
		protoUser.Height = *user.Height
	}
	if user.Weight != nil {
		protoUser.Weight = *user.Weight
	}
	if user.Budget != nil {
		protoUser.Budget = *user.Budget
	}
	if user.BJU != nil {
		protoUser.Bju = &proto_models.BJUModel{
			Protein: user.BJU.Protein,
			Fat:     user.BJU.Fat,
			Carbs:   user.BJU.Carbs,
		}
	}

	return protoUser
}
