package profile_management_storage

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/Android12349/food_recomendation/profile_managment_service/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func (s *ProfileManagementStorage) CreateUser(ctx context.Context, user *models.User) error {
	var bjuJSON []byte
	if user.BJU != nil {
		var err error
		bjuJSON, err = json.Marshal(user.BJU)
		if err != nil {
			return errors.Wrap(err, "marshal bju")
		}
	}

	query := squirrel.Insert(usersTableName).
		Columns(usersUsernameColumn, usersPasswordHashColumn, usersHeightColumn,
			usersWeightColumn, usersBJUColumn, usersBudgetColumn, usersPreferencesColumn).
		Values(user.Username, user.PasswordHash, user.Height, user.Weight,
			bjuJSON, user.Budget, user.Preferences).
		Suffix("RETURNING " + usersIDColumn + ", " + usersCreatedAtColumn).
		PlaceholderFormat(squirrel.Dollar)

	queryText, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "generate query error")
	}

	shard := s.getShardByUsername(user.Username)
	var createdAt sql.NullTime
	err = shard.QueryRow(ctx, queryText, args...).Scan(&user.ID, &createdAt)
	if err != nil {
		return errors.Wrap(err, "exec query error")
	}

	if createdAt.Valid {
		user.CreatedAt = createdAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	return nil
}

func (s *ProfileManagementStorage) GetUserByID(ctx context.Context, id int32) (*models.User, error) {
	query := squirrel.Select(usersIDColumn, usersUsernameColumn, usersPasswordHashColumn,
		usersHeightColumn, usersWeightColumn, usersBJUColumn, usersBudgetColumn,
		usersPreferencesColumn, usersCreatedAtColumn).
		From(usersTableName).
		Where(squirrel.Eq{usersIDColumn: id}).
		PlaceholderFormat(squirrel.Dollar)

	queryText, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "generate query error")
	}

	shard := s.getShard(id)
	var user models.User
	var bjuJSON []byte
	var height, weight, budget sql.NullInt32
	var createdAt sql.NullTime
	var found bool

	err = shard.QueryRow(ctx, queryText, args...).Scan(
		&user.ID, &user.Username, &user.PasswordHash,
		&height, &weight, &bjuJSON, &budget,
		&user.Preferences, &createdAt,
	)
	if err == nil {
		found = true
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.Wrap(err, "scan row error")
	}

	if !found {
		for _, shard := range s.shards {
			err = shard.QueryRow(ctx, queryText, args...).Scan(
				&user.ID, &user.Username, &user.PasswordHash,
				&height, &weight, &bjuJSON, &budget,
				&user.Preferences, &createdAt,
			)
			if err == nil {
				found = true
				break
			}
			if !errors.Is(err, pgx.ErrNoRows) {
				return nil, errors.Wrap(err, "scan row error")
			}
		}
	}

	if !found {
		return nil, errors.New("user not found")
	}

	if height.Valid {
		user.Height = &height.Int32
	}
	if weight.Valid {
		user.Weight = &weight.Int32
	}
	if budget.Valid {
		user.Budget = &budget.Int32
	}

	if len(bjuJSON) > 0 {
		var bju models.BJU
		err = json.Unmarshal(bjuJSON, &bju)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal bju")
		}
		user.BJU = &bju
	}

	if createdAt.Valid {
		user.CreatedAt = createdAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	return &user, nil
}

func (s *ProfileManagementStorage) UpdateUser(ctx context.Context, user *models.User) error {
	var bjuJSON []byte
	if user.BJU != nil {
		var err error
		bjuJSON, err = json.Marshal(user.BJU)
		if err != nil {
			return errors.Wrap(err, "marshal bju")
		}
	}

	query := squirrel.Update(usersTableName).
		Set(usersUsernameColumn, user.Username).
		Set(usersHeightColumn, user.Height).
		Set(usersWeightColumn, user.Weight).
		Set(usersBJUColumn, bjuJSON).
		Set(usersBudgetColumn, user.Budget).
		Set(usersPreferencesColumn, user.Preferences).
		Where(squirrel.Eq{usersIDColumn: user.ID}).
		PlaceholderFormat(squirrel.Dollar)

	queryText, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "generate query error")
	}

	shard := s.getShard(user.ID)
	result, err := shard.Exec(ctx, queryText, args...)
	if err != nil {
		return errors.Wrap(err, "exec query error")
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		for _, shard := range s.shards {
			result, err = shard.Exec(ctx, queryText, args...)
			if err != nil {
				return errors.Wrap(err, "exec query error")
			}
			if result.RowsAffected() > 0 {
				return nil
			}
		}
		return errors.New("user not found")
	}

	return nil
}

func (s *ProfileManagementStorage) DeleteUser(ctx context.Context, id int32) error {
	query := squirrel.Delete(usersTableName).
		Where(squirrel.Eq{usersIDColumn: id}).
		PlaceholderFormat(squirrel.Dollar)

	queryText, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "generate query error")
	}

	shard := s.getShard(id)
	result, err := shard.Exec(ctx, queryText, args...)
	if err != nil {
		return errors.Wrap(err, "exec query error")
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		// Пропускаем первый шард, так как он уже проверен
		for i := 1; i < len(s.shards); i++ {
			result, err = s.shards[i].Exec(ctx, queryText, args...)
			if err != nil {
				return errors.Wrap(err, "exec query error")
			}
			if result.RowsAffected() > 0 {
				return nil
			}
		}
		return errors.New("user not found")
	}

	return nil
}
