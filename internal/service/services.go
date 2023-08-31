package service

import (
	"avito-internship/internal/entity"
	"avito-internship/internal/repository"
	"avito-internship/internal/webapi"
	"context"
)

// Segment методы сервиса сегментов
type Segment interface {
	// CreateSegment метод, создающий сегмент,
	// на вход принимает название сегмента и [опционально] необходимый процент пользователей,
	// возвращает ошибку или nil
	CreateSegment(ctx context.Context, req entity.SegmentRequest) error

	// DeleteSegment метод, удаляющий сегмент,
	// на вход принимает название сегмента,
	// возвращает ошибку или nil
	DeleteSegment(ctx context.Context, req entity.SegmentRequest) error
}

// User методы сервиса пользователей
type User interface {
	// AddSegment метод, добавляющий пользователя в сегменты,
	// на вход принимает id пользователя, массив из названий сегментов и время нахождения пользователя в указанных сегментах,
	// возвращает ошибку или nil.
	// При отсутствии ttl, конечное время не указывается.
	// ttl задаётся в часах.
	AddSegment(ctx context.Context, req entity.UserAddToSegmentRequest) error

	// RemoveSegment метод, исключающий пользователя из сегментов,
	// на вход принимает id пользователя и массив из названий сегментов,
	// возвращает ошибку или nil.
	RemoveSegment(ctx context.Context, req entity.UserRemoveFromSegmentRequest) error

	// GetActiveSegments метод, возвращающий массив сегментов в которых состоит пользователь,
	// на вход принимает id пользователя и массив из названий сегментов,
	// помимо массива активных сегментов пользователя возвращает ошибку или nil.
	GetActiveSegments(ctx context.Context, req entity.UserActiveSegmentRequest) ([]string, error)
}

// Report методы сервиса отчетов
type Report interface {
	// GetUserHistory метод, составляющий историю операций за конкретный период времени,
	// на вход принимает месяц и год (int),
	// возвращает массив из полей отчета и их значений, также возвращает ошибку или nil.
	GetUserHistory(ctx context.Context, req entity.ReportRequest) ([]entity.ReportUserHistory, error)

	// MakeReportLink метод, создающий ссылку с отчетом в формате csv на Google Drive,
	// на вход принимает месяц и год (int),
	// возвращает ссылку на отчет в Google Drive и ошибку или nil.
	MakeReportLink(ctx context.Context, req entity.ReportRequest) (string, error)

	// MakeReportFile метод, создающий отчет в формате массива байтов,
	// на вход принимает месяц и год (int),
	// возвращает массив байтов и ошибку или nil.
	MakeReportFile(ctx context.Context, req entity.ReportRequest) ([]byte, error)
}

type Services struct {
	Segment Segment
	User    User
	Report  Report
}

type ServicesDependencies struct {
	Repos  *repository.Repositories
	GDrive webapi.GDrive
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Segment: NewSegmentService(deps.Repos.SegmentRepo),
		User:    NewUserService(deps.Repos.UserRepo),
		Report:  NewReportService(deps.Repos.ReportRepo, deps.GDrive),
	}
}
