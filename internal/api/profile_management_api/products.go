package profile_management_api

import (
	"context"
	"log"

	"github.com/Android12349/food_recomendation/profile_managment_service/internal/models"
	proto_models "github.com/Android12349/food_recomendation/profile_managment_service/internal/pb/models"
	"github.com/Android12349/food_recomendation/profile_managment_service/internal/pb/profile_management_api"
	"github.com/samber/lo"
)

func (s *ProfileManagementAPI) CreateProduct(ctx context.Context, req *profile_management_api.CreateProductRequest) (*profile_management_api.CreateProductResponse, error) {
	log.Printf("Received CreateProduct request for user_id: %d", req.Product.UserId)

	product := mapProductCreateModelToModel(req.Product)

	err := s.profileService.CreateProduct(ctx, product)
	if err != nil {
		return &profile_management_api.CreateProductResponse{}, err
	}

	return &profile_management_api.CreateProductResponse{
		Product: mapProductModelToProto(product),
	}, nil
}

func (s *ProfileManagementAPI) GetProducts(ctx context.Context, req *profile_management_api.GetProductsRequest) (*profile_management_api.GetProductsResponse, error) {
	log.Printf("Received GetProducts request for user_id: %d", req.UserId)

	products, err := s.profileService.GetProductsByUserID(ctx, req.UserId)
	if err != nil {
		return &profile_management_api.GetProductsResponse{}, err
	}

	return &profile_management_api.GetProductsResponse{
		Products: lo.Map(products, func(p *models.Product, _ int) *proto_models.ProductModel {
			return mapProductModelToProto(p)
		}),
	}, nil
}

func (s *ProfileManagementAPI) UpdateProduct(ctx context.Context, req *profile_management_api.UpdateProductRequest) (*profile_management_api.UpdateProductResponse, error) {
	log.Printf("Received UpdateProduct request for ID: %d", req.Id)

	existingProduct, err := s.profileService.GetProductByID(ctx, req.Id)
	if err != nil {
		return &profile_management_api.UpdateProductResponse{}, err
	}

	product := mapProductUpdateModelToModel(req.Product, req.Id, existingProduct.UserID)

	err = s.profileService.UpdateProduct(ctx, product)
	if err != nil {
		return &profile_management_api.UpdateProductResponse{}, err
	}

	updatedProduct, err := s.profileService.GetProductByID(ctx, req.Id)
	if err != nil {
		return &profile_management_api.UpdateProductResponse{}, err
	}

	return &profile_management_api.UpdateProductResponse{
		Product: mapProductModelToProto(updatedProduct),
	}, nil
}

func (s *ProfileManagementAPI) DeleteProduct(ctx context.Context, req *profile_management_api.DeleteProductRequest) (*profile_management_api.DeleteProductResponse, error) {
	log.Printf("Received DeleteProduct request for ID: %d", req.Id)

	err := s.profileService.DeleteProduct(ctx, req.Id)
	if err != nil {
		return &profile_management_api.DeleteProductResponse{}, err
	}

	return &profile_management_api.DeleteProductResponse{}, nil
}

func mapProductCreateModelToModel(protoProduct *proto_models.ProductCreateModel) *models.Product {
	product := &models.Product{
		UserID: protoProduct.UserId,
		Name:   protoProduct.Name,
	}

	if protoProduct.Calories != 0 {
		calories := protoProduct.Calories
		product.Calories = &calories
	}
	if protoProduct.Protein != 0 {
		protein := protoProduct.Protein
		product.Protein = &protein
	}
	if protoProduct.Fat != 0 {
		fat := protoProduct.Fat
		product.Fat = &fat
	}
	if protoProduct.Carbs != 0 {
		carbs := protoProduct.Carbs
		product.Carbs = &carbs
	}

	return product
}

func mapProductUpdateModelToModel(protoProduct *proto_models.ProductUpdateModel, id int32, userID int32) *models.Product {
	product := &models.Product{
		ID:     id,
		UserID: userID,
		Name:   protoProduct.Name,
	}

	if protoProduct.Calories != 0 {
		calories := protoProduct.Calories
		product.Calories = &calories
	}
	if protoProduct.Protein != 0 {
		protein := protoProduct.Protein
		product.Protein = &protein
	}
	if protoProduct.Fat != 0 {
		fat := protoProduct.Fat
		product.Fat = &fat
	}
	if protoProduct.Carbs != 0 {
		carbs := protoProduct.Carbs
		product.Carbs = &carbs
	}

	return product
}

func mapProductModelToProto(product *models.Product) *proto_models.ProductModel {
	protoProduct := &proto_models.ProductModel{
		Id:        product.ID,
		UserId:    product.UserID,
		Name:      product.Name,
		CreatedAt: product.CreatedAt,
	}

	if product.Calories != nil {
		protoProduct.Calories = *product.Calories
	}
	if product.Protein != nil {
		protoProduct.Protein = *product.Protein
	}
	if product.Fat != nil {
		protoProduct.Fat = *product.Fat
	}
	if product.Carbs != nil {
		protoProduct.Carbs = *product.Carbs
	}

	return protoProduct
}
