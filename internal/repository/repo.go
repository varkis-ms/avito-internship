package repository

import (
	"avito-internship/internal/entity"
	"avito-internship/internal/repository/pgdb"
	"avito-internship/pkg/database/postgresdb"
	"context"
)

// SegmentRepo Методы репозитория сегментов
type SegmentRepo interface {
	// CreateSegment метод создания сегмента, на вход принимает название,
	// возвращает id созданного сегмента и ошибку бд или nil
	CreateSegment(ctx context.Context, segment string) (int, error)

	// DeleteSegment метод удаления сегмента, на вход принимает название,
	// возвращает ошибку бд или nil
	DeleteSegment(ctx context.Context, segment string) error

	// CheckExistSegment метод проверки сегмента на существование, на вход принимает название,
	// возвращает true если сегмент существует, иначе false, и ошибку бд или nil
	CheckExistSegment(ctx context.Context, segment string) (bool, error)

	// RandomUserToSegment метод добавления случайных N% пользователей в сегмент,
	// на вход принимает id сегмента и необходимый процент пользователей,
	// возвращает ошибку бд или nil
	RandomUserToSegment(ctx context.Context, segmentId int, percent float32) error
}

// UserRepo Методы репозитория пользователей
type UserRepo interface {
	// AddSegmentToUser метод добавления пользователя в сегменты,
	// на вход принимает id пользователя, массив из id сегментов и время нахождения пользователя в указанных сегментах,
	// возвращает ошибку бд или nil.
	// При отсутствии ttl, конечное время не указывается.
	// ttl задаётся в часах.
	AddSegmentToUser(ctx context.Context, id int, segments []int, ttl int) error

	// RemoveSegmentFromUser метод исключения пользователя из сегментов,
	// на вход принимает id пользователя и массив из id сегментов,
	// возвращает ошибку бд или nil.
	RemoveSegmentFromUser(ctx context.Context, id int, segments []int) error

	// GetActiveSegmentsIdByName метод получения активных сегментов сервиса,
	// на вход принимает массив из id сегментов,
	// возвращает массив из id сегментов и ошибку бд или nil.
	GetActiveSegmentsIdByName(ctx context.Context, segments []string) ([]int, error)

	// GetActiveSegmentFromUser метод получения активных сегментов пользователя,
	// на вход принимает id пользователя,
	// возвращает массив из названий сегментов и ошибку бд или nil.
	GetActiveSegmentFromUser(ctx context.Context, id int) ([]string, error)

	// CheckExistUser метод проверки существования пользователя,
	// на вход принимает id пользователя,
	// возвращает ошибку бд (в том числе и при не существовании пользователя) или nil.
	CheckExistUser(ctx context.Context, id int) error
}

// ReportRepo Методы репозитория отчета
type ReportRepo interface {
	// GetSegmentHistoryFromUser метод получения истории пользователей (вхождение/исключение из сегментов),
	// на вход принимает месяц и год,
	// возвращает массив из ReportUserHistory ошибку бд или nil.
	GetSegmentHistoryFromUser(ctx context.Context, month int, year int) ([]entity.ReportUserHistory, error)
}

type Repositories struct {
	SegmentRepo
	UserRepo
	ReportRepo
}

func NewRepositories(pg *postgresdb.Postgres) *Repositories {
	return &Repositories{
		SegmentRepo: pgdb.NewSegmentRepo(pg),
		UserRepo:    pgdb.NewUserRepo(pg),
		ReportRepo:  pgdb.NewReportRepo(pg),
	}
}
