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

func TestCreateSegment(t *testing.T) {
	type args struct {
		ctx     context.Context
		segment string
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         int
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{ctx: context.Background(),
				segment: "Test_Segment",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(1)
				m.ExpectQuery("INSERT").
					WithArgs(args.segment).WillReturnRows(rows)
			},
			wantErr: false,
			want:    1,
		},
		{
			name: "Segment_exist",
			args: args{ctx: context.Background(),
				segment: "Test_Segment",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("INSERT").
					WithArgs(args.segment).WillReturnError(pgx.ErrNoRows)
			},
			wantErr: false,
			want:    0,
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
			segmentRepoMock := pgdb.NewSegmentRepo(postgresMock)
			got, err := segmentRepoMock.CreateSegment(tc.args.ctx, tc.args.segment)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDeleteSegment(t *testing.T) {
	type args struct {
		ctx     context.Context
		segment string
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
				segment: "Test_Segment",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectBegin()

				rows := pgxmock.NewRows([]string{"id"}).AddRow(1)
				m.ExpectQuery("UPDATE").
					WithArgs("now()", args.segment).WillReturnRows(rows)

				m.ExpectExec("UPDATE").
					WithArgs("now()", 1, "now()").
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

				m.ExpectCommit()

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
			segmentRepoMock := pgdb.NewSegmentRepo(postgresMock)
			err := segmentRepoMock.DeleteSegment(tc.args.ctx, tc.args.segment)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckExistSegment(t *testing.T) {
	type args struct {
		ctx     context.Context
		segment string
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         bool
		wantErr      bool
	}{
		{
			name: "OK_false",
			args: args{ctx: context.Background(),
				segment: "Test_Segment",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"exists"}).AddRow(false)
				m.ExpectQuery("SELECT").
					WithArgs(args.segment, "now()").WillReturnRows(rows)
			},
			wantErr: false,
			want:    false,
		},
		{
			name: "OK_true",
			args: args{ctx: context.Background(),
				segment: "Test_Segment",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.QueryRow(args.ctx, "INSERT INTO segments (name) VALUES ($1)", args.segment)

				rows := pgxmock.NewRows([]string{"exists"}).AddRow(true)
				m.ExpectQuery("SELECT").
					WithArgs(args.segment, "now()").WillReturnRows(rows)
			},
			wantErr: false,
			want:    true,
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
			segmentRepoMock := pgdb.NewSegmentRepo(postgresMock)
			got, err := segmentRepoMock.CheckExistSegment(tc.args.ctx, tc.args.segment)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestRandomUserToSegment(t *testing.T) {
	type args struct {
		ctx       context.Context
		segmentId int
		percent   float32
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
				segmentId: 1,
				percent:   0.5,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewResult("INSERT", 0)
				m.ExpectExec("INSERT").
					WithArgs(args.percent).WillReturnResult(rows)
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
			segmentRepoMock := pgdb.NewSegmentRepo(postgresMock)
			err := segmentRepoMock.RandomUserToSegment(tc.args.ctx, tc.args.segmentId, tc.args.percent)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
