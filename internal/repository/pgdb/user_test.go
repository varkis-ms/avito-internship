package pgdb_test

import (
	"avito-internship/internal/repository/pgdb"
	"avito-internship/pkg/database/postgresdb"
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoveSegmentFromUser(t *testing.T) {
	type args struct {
		ctx      context.Context
		id       int
		segments []int
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{ctx: context.Background(),
				id:       1,
				segments: []int{1, 2},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec("UPDATE").
					WithArgs("now()", "now()", args.id, args.segments[0], args.segments[1]).
					WillReturnResult(pgxmock.NewResult("UPDATE", 2))
			},
			wantErr: false,
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
			userRepoMock := pgdb.NewUserRepo(postgresMock)
			err := userRepoMock.RemoveSegmentFromUser(tc.args.ctx, tc.args.id, tc.args.segments)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetActiveSegmentFromUser(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantErr      bool
		want         []string
	}{
		{
			name: "OK",
			args: args{ctx: context.Background(),
				id: 1,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"name"}).AddRow("test_segment_1").AddRow("test_segment_2")
				m.ExpectQuery("SELECT").
					WithArgs("now()", args.id).WillReturnRows(rows)
			},
			wantErr: false,
			want:    []string{"test_segment_1", "test_segment_2"},
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
			userRepoMock := pgdb.NewUserRepo(postgresMock)
			got, err := userRepoMock.GetActiveSegmentFromUser(tc.args.ctx, tc.args.id)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCheckExistUser(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{ctx: context.Background(),
				id: 1,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"exists"}).AddRow(true)
				m.ExpectQuery("SELECT").
					WithArgs(args.id).WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "Error_does_not_exist)",
			args: args{ctx: context.Background(),
				id: 1,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("SELECT").
					WithArgs(args.id).
					WillReturnError(pgx.ErrNoRows)
			},
			wantErr: true,
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
			userRepoMock := pgdb.NewUserRepo(postgresMock)
			err := userRepoMock.CheckExistUser(tc.args.ctx, tc.args.id)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetActiveSegmentsIdByName(t *testing.T) {
	type args struct {
		ctx      context.Context
		segments []string
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantErr      bool
		want         []int
	}{
		{
			name: "OK",
			args: args{ctx: context.Background(),
				segments: []string{"test_segment_1"},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(1)
				m.ExpectQuery("SELECT").
					WithArgs(args.segments[0], "now()").
					WillReturnRows(rows)
			},
			wantErr: false,
			want:    []int{1},
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
			userRepoMock := pgdb.NewUserRepo(postgresMock)
			got, err := userRepoMock.GetActiveSegmentsIdByName(tc.args.ctx, tc.args.segments)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.want, got)
		})
	}
}
