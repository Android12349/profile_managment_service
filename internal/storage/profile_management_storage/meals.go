package profile_management_storage

import (
	"context"
	"database/sql"

	"github.com/Android12349/food_recomendation/profile_managment_service/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func (s *ProfileManagementStorage) CreateMeal(ctx context.Context, meal *models.Meal) error {
	query := squirrel.Insert(mealsTableName).
		Columns(mealsUserIDColumn, mealsNameColumn, mealsProductIDsColumn).
		Values(meal.UserID, meal.Name, meal.ProductIDs).
		Suffix("RETURNING " + mealsIDColumn + ", " + mealsCreatedAtColumn).
		PlaceholderFormat(squirrel.Dollar)

	queryText, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "generate query error")
	}

	shard := s.getShard(meal.UserID)
	var createdAt sql.NullTime
	err = shard.QueryRow(ctx, queryText, args...).Scan(&meal.ID, &createdAt)
	if err != nil {
		return errors.Wrap(err, "exec query error")
	}

	if createdAt.Valid {
		meal.CreatedAt = createdAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	return nil
}

func (s *ProfileManagementStorage) GetMealsByUserID(ctx context.Context, userID int32) ([]*models.Meal, error) {
	query := squirrel.Select(mealsIDColumn, mealsUserIDColumn, mealsNameColumn,
		mealsProductIDsColumn, mealsCreatedAtColumn).
		From(mealsTableName).
		Where(squirrel.Eq{mealsUserIDColumn: userID}).
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

	var meals []*models.Meal
	for rows.Next() {
		var meal models.Meal
		var productIDs []int32
		var createdAt sql.NullTime

		err := rows.Scan(
			&meal.ID, &meal.UserID, &meal.Name,
			&productIDs, &createdAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "scan row error")
		}

		meal.ProductIDs = productIDs

		if createdAt.Valid {
			meal.CreatedAt = createdAt.Time.Format("2006-01-02T15:04:05Z07:00")
		}

		meals = append(meals, &meal)
	}

	return meals, nil
}

func (s *ProfileManagementStorage) GetMealByID(ctx context.Context, id int32) (*models.Meal, error) {
	query := squirrel.Select(mealsIDColumn, mealsUserIDColumn, mealsNameColumn,
		mealsProductIDsColumn, mealsCreatedAtColumn).
		From(mealsTableName).
		Where(squirrel.Eq{mealsIDColumn: id}).
		PlaceholderFormat(squirrel.Dollar)

	queryText, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "generate query error")
	}

	var meal models.Meal
	var productIDs []int32
	var createdAt sql.NullTime

	var found bool
	for _, shard := range s.shards {
		err = shard.QueryRow(ctx, queryText, args...).Scan(
			&meal.ID, &meal.UserID, &meal.Name,
			&productIDs, &createdAt,
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
		return nil, errors.New("meal not found")
	}

	meal.ProductIDs = productIDs

	if createdAt.Valid {
		meal.CreatedAt = createdAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	return &meal, nil
}

func (s *ProfileManagementStorage) UpdateMeal(ctx context.Context, meal *models.Meal) error {
	query := squirrel.Update(mealsTableName).
		Set(mealsNameColumn, meal.Name).
		Set(mealsProductIDsColumn, meal.ProductIDs).
		Where(squirrel.Eq{mealsIDColumn: meal.ID}).
		PlaceholderFormat(squirrel.Dollar)

	queryText, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "generate query error")
	}

	tempMeal, err := s.GetMealByID(ctx, meal.ID)
	if err != nil {
		return err
	}
	shard := s.getShard(tempMeal.UserID)
	_, err = shard.Exec(ctx, queryText, args...)
	if err != nil {
		return errors.Wrap(err, "exec query error")
	}

	return nil
}

func (s *ProfileManagementStorage) DeleteMeal(ctx context.Context, id int32) error {
	query := squirrel.Delete(mealsTableName).
		Where(squirrel.Eq{mealsIDColumn: id}).
		PlaceholderFormat(squirrel.Dollar)

	queryText, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "generate query error")
	}

	tempMeal, err := s.GetMealByID(ctx, id)
	if err != nil {
		return err
	}
	shard := s.getShard(tempMeal.UserID)
	_, err = shard.Exec(ctx, queryText, args...)
	if err != nil {
		return errors.Wrap(err, "exec query error")
	}

	return nil
}
