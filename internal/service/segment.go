package service

import (
	"avito-internship/internal/entity"
	"avito-internship/internal/repository"
	"context"
	"fmt"
)

type SegmentService struct {
	segmentRepo repository.SegmentRepo
}

func NewSegmentService(segmentRepo repository.SegmentRepo) *SegmentService {
	return &SegmentService{segmentRepo: segmentRepo}
}

func (s *SegmentService) CreateSegment(ctx context.Context, req entity.SegmentRequest) error {
	segmentId, err := s.segmentRepo.CreateSegment(ctx, req.Segment)
	if err != nil {
		return fmt.Errorf("segmentRepo.CreateSegment: %w", err)
	}

	if segmentId != 0 && req.Percent != 0.0 && (req.Percent < 1.0 || req.Percent > 0.0) {
		err = s.segmentRepo.RandomUserToSegment(ctx, segmentId, req.Percent)
		if err != nil {
			return fmt.Errorf("segmentRepo.RandomUserToSegment: %w", err)
		}
	}

	return nil
}

func (s *SegmentService) DeleteSegment(ctx context.Context, req entity.SegmentRequest) error {
	exist, err := s.segmentRepo.CheckExistSegment(ctx, req.Segment)
	if err != nil {
		return fmt.Errorf("segmentRepo.CheckExistSegment: %w", err)
	}

	if !exist {
		return nil
	}

	err = s.segmentRepo.DeleteSegment(ctx, req.Segment)
	if err != nil {
		return fmt.Errorf("segmentRepo.DeleteSegment: %w", err)
	}

	return nil
}
