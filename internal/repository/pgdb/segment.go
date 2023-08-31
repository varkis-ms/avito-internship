package pgdb

import (
	"avito-internship/pkg/database/postgresdb"
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type SegmentRepo struct {
	*postgresdb.Postgres
}

func NewSegmentRepo(pg *postgresdb.Postgres) *SegmentRepo {
	return &SegmentRepo{pg}
}

func (r *SegmentRepo) CreateSegment(ctx context.Context, segment string) (int, error) {
	sql, args, _ := r.Builder.
		Insert("segments").
		Columns("name").
		Values(segment).
		Suffix("ON CONFLICT DO NOTHING").
		Suffix("RETURNING id").
		ToSql()

	var segmentId int
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&segmentId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, err
	}

	return segmentId, nil
}

func (r *SegmentRepo) DeleteSegment(ctx context.Context, segment string) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	sql, args, _ := r.Builder.
		Update("segments").
		Set("deleted_at", "now()").
		Where("name = ?", segment).
		Suffix("RETURNING id").
		ToSql()

	var segmentId int
	err = tx.QueryRow(ctx, sql, args...).Scan(&segmentId)
	if err != nil {
		return err
	}

	sql, args, _ = r.Builder.
		Update("users_segment").
		Set("left_at", "now()").
		Where("segment_id = ?", segmentId).
		Where(sq.Or{
			sq.Eq{"left_at": nil},
			sq.Gt{"left_at": "now()"},
		}).
		ToSql()

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

func (r *SegmentRepo) CheckExistSegment(ctx context.Context, segment string) (bool, error) {
	sql, args, _ := r.Builder.
		Select("1").
		Prefix("SELECT EXISTS (").
		From("segments").
		Where("name = ?", segment).
		Where(sq.Or{
			sq.Eq{"deleted_at": nil},
			sq.Gt{"deleted_at": "now()"},
		}).
		Suffix(")").
		ToSql()

	var exist bool
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&exist)

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return false, err
		}
	}

	return exist, nil
}

func (r *SegmentRepo) RandomUserToSegment(ctx context.Context, segmentId int, percent float32) error {
	subQuery := fmt.Sprintf("id, %v", segmentId)

	sql, args, _ := r.Builder.
		Insert("users_segment").
		Columns("user_id", "segment_id").
		Select(
			sq.Select(subQuery).
				From("users").
				Where(sq.Expr("random() < ?", percent))).
		ToSql()

	_, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
