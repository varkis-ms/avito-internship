package service

import (
	"avito-internship/internal/apperror"
	"avito-internship/internal/entity"
	"avito-internship/internal/repository"
	"avito-internship/internal/webapi"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"sort"
)

type ReportService struct {
	reportRepo repository.ReportRepo
	gDrive     webapi.GDrive
}

func NewReportService(reportRepo repository.ReportRepo, gDrive webapi.GDrive) *ReportService {
	return &ReportService{
		reportRepo: reportRepo,
		gDrive:     gDrive,
	}
}

func (s *ReportService) GetUserHistory(ctx context.Context, req entity.ReportRequest) ([]entity.ReportUserHistory, error) {
	userHistory, err := s.reportRepo.GetSegmentHistoryFromUser(ctx, req.Month, req.Year)
	if err != nil {
		return nil, fmt.Errorf("reportRepo.GetSegmentHistoryFromUser: %w", err)
	}
	sortByDate := func(i, j int) bool {
		return userHistory[i].Date.Before(userHistory[j].Date)
	}

	sort.SliceStable(userHistory, sortByDate)

	return userHistory, nil
}

func (s *ReportService) MakeReportLink(ctx context.Context, req entity.ReportRequest) (string, error) {
	if !s.gDrive.IsAvailable() {
		return "", apperror.ErrGDriveNotAvailable
	}

	file, err := s.MakeReportFile(ctx, req)
	if err != nil {
		return "", fmt.Errorf("reportRepo.MakeReportFile: %w", err)
	}

	url, err := s.gDrive.UploadCSVFile(ctx, fmt.Sprintf("report_%d_%d.csv", req.Month, req.Year), file)
	if err != nil {
		return "", fmt.Errorf("gDrive.UploadCSVFile: %w", err)
	}

	return url, nil
}

func (s *ReportService) MakeReportFile(ctx context.Context, req entity.ReportRequest) ([]byte, error) {
	report, err := s.GetUserHistory(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("reportService.GetUserHistory: %w", err)
	}

	b := bytes.Buffer{}
	w := csv.NewWriter(&b)

	err = w.Write([]string{
		"user_id",
		"segment",
		"operation",
		"date",
	})
	if err != nil {
		return nil, fmt.Errorf("reportService.MakeReportFile - w.Write Header: %w", err)
	}

	for _, item := range report {
		err = w.Write([]string{
			item.UserId,
			item.Segment,
			item.Operation,
			item.Date.String(),
		})
		if err != nil {
			return nil, fmt.Errorf("reportService.MakeReportFile - w.Write: %w", err)
		}
	}

	w.Flush()
	if err = w.Error(); err != nil {
		return nil, fmt.Errorf("reportService.MakeReportFile - w.Error(): %w", err)
	}

	return b.Bytes(), nil
}
