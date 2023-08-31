package pgdb

import (
	"avito-internship/internal/apperror"
	"avito-internship/internal/utils"
	"avito-internship/pkg/database/postgresdb"
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type UserRepo struct {
	*postgresdb.Postgres
}

func NewUserRepo(pg *postgresdb.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) AddSegmentToUser(ctx context.Context, id int, segments []int, ttl int) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	sql, args, _ := r.Builder.
		Insert("users").
		Columns("id").
		Values(id).
		Suffix("ON CONFLICT DO NOTHING").
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	}

	sql, args, _ = r.Builder.
		Select("segment_id").
		From("users_segment").
		Where("user_id = ?", id).
		Where(sq.Eq{"segment_id": segments}).
		Where(sq.Or{
			sq.Eq{"left_at": nil},
			sq.Gt{"left_at": "now()"},
		}).
		ToSql()

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var segmentID int
		err = rows.Scan(&segmentID)
		if err != nil {
			return err
		}
		segments = append(segments, segmentID)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	idToInsert := utils.UniqueValues(segments)

	sqlQuery := r.Builder.Insert("users_segment")
	if ttl > 0 {
		ttlInSql := sq.Expr(fmt.Sprintf("now() + INTERVAL '%d hours'", ttl))
		sqlQuery = sqlQuery.Columns("user_id", "segment_id", "left_at")
		for _, segmentId := range idToInsert {
			sqlQuery = sqlQuery.Values(id, segmentId, ttlInSql)
		}
	} else {
		sqlQuery = sqlQuery.Columns("user_id", "segment_id")
		for _, segmentId := range idToInsert {
			sqlQuery = sqlQuery.Values(id, segmentId)
		}
	}

	sql, args, _ = sqlQuery.ToSql()
	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) RemoveSegmentFromUser(ctx context.Context, id int, segments []int) error {
	sql, args, _ := r.Builder.
		Update("users_segment").
		Set("left_at", "now()").
		Where(sq.Or{
			sq.Eq{"left_at": nil},
			sq.Gt{"left_at": "now()"},
		}).
		Where("user_id = ?", id).
		Where(sq.Eq{"segment_id": segments}).
		ToSql()

	_, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) GetActiveSegmentFromUser(ctx context.Context, id int) ([]string, error) {
	sql, args, _ := r.Builder.
		Select("s.name").
		From("segments AS s").
		Join("users_segment AS us ON s.id = us.segment_id").
		Where(sq.Or{
			sq.Eq{"us.left_at": nil},
			sq.Gt{"us.left_at": "now()"},
		}).
		Where(sq.Eq{"us.user_id": id}).
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var segmentNames []string
	for rows.Next() {
		var segment string
		err = rows.Scan(&segment)
		if err != nil {
			return nil, err
		}
		segmentNames = append(segmentNames, segment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return segmentNames, nil
}

func (r *UserRepo) CheckExistUser(ctx context.Context, id int) error {
	sql, args, _ := r.Builder.
		Select("1").
		Prefix("SELECT EXISTS (").
		From("users").
		Where("id = ?", id).
		Suffix(")").
		ToSql()

	var exist bool
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&exist)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.ErrNoUser
		}

		return err
	}

	return nil
}

func (r *UserRepo) GetActiveSegmentsIdByName(ctx context.Context, segments []string) ([]int, error) {
	sql, args, _ := r.Builder.
		Select("id").
		From("segments").
		Where(sq.Eq{"name": segments}).
		Where(sq.Or{
			sq.Eq{"deleted_at": nil},
			sq.Gt{"deleted_at": "now()"},
		}).
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var segmentsIds []int
	for rows.Next() {
		var segmentId int
		err = rows.Scan(&segmentId)
		if err != nil {
			return nil, err
		}
		segmentsIds = append(segmentsIds, segmentId)
	}

	return segmentsIds, nil
}
