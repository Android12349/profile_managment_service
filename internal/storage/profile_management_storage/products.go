package profile_management_storage

import (
	"context"
	"database/sql"

	"github.com/Android12349/food_recomendation/profile_managment_service/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func (s *ProfileManagementStorage) CreateProduct(ctx context.Context, product *models.Product) error {
	query := squirrel.Insert(productsTableName).
		Columns(productsUserIDColumn, productsNameColumn, productsCaloriesColumn,
			productsProteinColumn, productsFatColumn, productsCarbsColumn).
		Values(product.UserID, product.Name, product.Calories,
			product.Protein, product.Fat, product.Carbs).
		Suffix("RETURNING " + productsIDColumn + ", " + productsCreatedAtColumn).
		PlaceholderFormat(squirrel.Dollar)

	queryText, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "generate query error")
	}

	shard := s.getShard(product.UserID)
	var createdAt sql.NullTime
	err = shard.QueryRow(ctx, queryText, args...).Scan(&product.ID, &createdAt)
	if err != nil {
		return errors.Wrap(err, "exec query error")
	}

	if createdAt.Valid {
		product.CreatedAt = createdAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	return nil
}

func (s *ProfileManagementStorage) GetProductsByUserID(ctx context.Context, userID int32) ([]*models.Product, error) {
	query := squirrel.Select(productsIDColumn, productsUserIDColumn, productsNameColumn,
		productsCaloriesColumn, productsProteinColumn, productsFatColumn,
		productsCarbsColumn, productsCreatedAtColumn).
		From(productsTableName).
		Where(squirrel.Eq{productsUserIDColumn: userID}).
		PlaceholderFormat(squirrel.Dollar)

	queryText, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "generate query error")
	}

	shard := s.getShard(userID)
	rows, err := shard.Query(ctx, queryText, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var product models.Product
		var calories, protein, fat, carbs sql.NullInt32
		var createdAt sql.NullTime

		err := rows.Scan(
			&product.ID, &product.UserID, &product.Name,
			&calories, &protein, &fat,
			&carbs, &createdAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "scan row error")
		}

		if calories.Valid {
			product.Calories = &calories.Int32
		}
		if protein.Valid {
			product.Protein = &protein.Int32
		}
		if fat.Valid {
			product.Fat = &fat.Int32
		}
		if carbs.Valid {
			product.Carbs = &carbs.Int32
		}

		if createdAt.Valid {
			product.CreatedAt = createdAt.Time.Format("2006-01-02T15:04:05Z07:00")
		}

		products = append(products, &product)
	}

	return products, nil
}

func (s *ProfileManagementStorage) GetProductByID(ctx context.Context, id int32) (*models.Product, error) {
	query := squirrel.Select(productsIDColumn, productsUserIDColumn, productsNameColumn,
		productsCaloriesColumn, productsProteinColumn, productsFatColumn,
		productsCarbsColumn, productsCreatedAtColumn).
		From(productsTableName).
		Where(squirrel.Eq{productsIDColumn: id}).
		PlaceholderFormat(squirrel.Dollar)

	queryText, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "generate query error")
	}

	var product models.Product
	var calories, protein, fat, carbs sql.NullInt32
	var createdAt sql.NullTime
	var found bool

	for _, shard := range s.shards {
		err = shard.QueryRow(ctx, queryText, args...).Scan(
			&product.ID, &product.UserID, &product.Name,
			&calories, &protein, &fat,
			&carbs, &createdAt,
		)
		if err == nil {
			found = true
			break
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(err, "scan row error")
		}
	}

	if !found {
		return nil, errors.New("product not found")
	}

	if calories.Valid {
		product.Calories = &calories.Int32
	}
	if protein.Valid {
		product.Protein = &protein.Int32
	}
	if fat.Valid {
		product.Fat = &fat.Int32
	}
	if carbs.Valid {
		product.Carbs = &carbs.Int32
	}

	if createdAt.Valid {
		product.CreatedAt = createdAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	return &product, nil
}

func (s *ProfileManagementStorage) UpdateProduct(ctx context.Context, product *models.Product) error {
	query := squirrel.Update(productsTableName).
		Set(productsNameColumn, product.Name).
		Set(productsCaloriesColumn, product.Calories).
		Set(productsProteinColumn, product.Protein).
		Set(productsFatColumn, product.Fat).
		Set(productsCarbsColumn, product.Carbs).
		Where(squirrel.Eq{productsIDColumn: product.ID}).
		PlaceholderFormat(squirrel.Dollar)

	queryText, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "generate query error")
	}

	tempProduct, err := s.GetProductByID(ctx, product.ID)
	if err != nil {
		return err
	}
	shard := s.getShard(tempProduct.UserID)
	_, err = shard.Exec(ctx, queryText, args...)
	if err != nil {
		return errors.Wrap(err, "exec query error")
	}

	return nil
}

func (s *ProfileManagementStorage) DeleteProduct(ctx context.Context, id int32) error {
	query := squirrel.Delete(productsTableName).
		Where(squirrel.Eq{productsIDColumn: id}).
		PlaceholderFormat(squirrel.Dollar)

	queryText, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "generate query error")
	}

	tempProduct, err := s.GetProductByID(ctx, id)
	if err != nil {
		return err
	}
	shard := s.getShard(tempProduct.UserID)
	_, err = shard.Exec(ctx, queryText, args...)
	if err != nil {
		return errors.Wrap(err, "exec query error")
	}

	return nil
}
