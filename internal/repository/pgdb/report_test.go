package pgdb_test

import (
	"avito-internship/internal/entity"
	"avito-internship/internal/repository/pgdb"
	"avito-internship/pkg/database/postgresdb"
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSegmentHistoryFromUser(t *testing.T) {
	type args struct {
		ctx   context.Context
		month int
		year  int
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         []entity.ReportUserHistory
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{ctx: context.Background(),
				month: 9,
				year:  2023,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{
					"user_id", "added_at",
					"name", "left_at",
				})
				m.ExpectQuery("SELECT").
					WithArgs(args.month, args.year).
					WillReturnRows(rows)
			},
			wantErr: false,
			want:    nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tc.mockBehavior(poolMock, tc.args)

			postgresMock := &postgresdb.Postgres{
				Builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
				Pool:    poolMock,
			}
			reportRepoMock := pgdb.NewReportRepo(postgresMock)
			got, err := reportRepoMock.GetSegmentHistoryFromUser(tc.args.ctx, tc.args.month, tc.args.year)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.want, got)
		})
	}
}
