package pgdb

import (
	"avito-internship/internal/entity"
	"avito-internship/internal/utils"
	"avito-internship/pkg/database/postgresdb"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

const (
	operationAdd    = "add"
	operationRemove = "remove"
)

type ReportRepo struct {
	*postgresdb.Postgres
}

func NewReportRepo(pg *postgresdb.Postgres) *ReportRepo {
	return &ReportRepo{pg}
}

func (r *ReportRepo) GetSegmentHistoryFromUser(ctx context.Context, month int, year int) ([]entity.ReportUserHistory, error) {
	sql, args, _ := r.Builder.
		Select("us.user_id", "us.added_at", "s.name", "us.left_at").
		From("users_segment AS us").
		Join("segments AS s ON s.id = us.segment_id").
		Where(sq.Eq{"extract(month from us.added_at)": month}).
		Where(sq.Eq{"extract(year from us.added_at)": year}).
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err

	}
	defer rows.Close()

	pgColumns := rows.FieldDescriptions()
	values := make([]interface{}, len(pgColumns))
	for i := range pgColumns {
		values[i] = new(interface{})
	}

	var results []entity.ReportUserHistory

	for rows.Next() {
		if err = rows.Scan(values...); err != nil {
			return nil, err
		}

		historyAdd := entity.ReportUserHistory{}
		historyRemove := entity.ReportUserHistory{}
		removeExist := false
		for i, col := range pgColumns {
			val := *values[i].(*interface{})

			if val != nil {
				switch col.Name {
				case "added_at":
					timestamp, err := utils.PgTimestampConverter(val)
					if err != nil {
						return nil, err
					}
					historyAdd.Date = timestamp
					historyAdd.Operation = operationAdd
				case "user_id":
					historyAdd.UserId = fmt.Sprint(val)
				case "name":
					historyAdd.Segment = val.(string)
				case "left_at":
					timestamp, err := utils.PgTimestampConverter(val)
					if err != nil {
						return nil, err
					}
					historyRemove = historyAdd
					historyRemove.Date = timestamp
					historyRemove.Operation = operationRemove
					removeExist = true
				}
			}
		}
		results = append(results, historyAdd)
		if removeExist {
			results = append(results, historyRemove)
		}
	}

	return results, nil
}
