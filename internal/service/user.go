package service

import (
	"avito-internship/internal/apperror"
	"avito-internship/internal/entity"
	"avito-internship/internal/repository"
	"context"
	"fmt"
)

type UserService struct {
	userRepo repository.UserRepo
}

func NewUserService(userRepo repository.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) AddSegment(ctx context.Context, req entity.UserAddToSegmentRequest) error {
	segmentsId, err := s.userRepo.GetActiveSegmentsIdByName(ctx, req.Segments)
	if err != nil {
		return fmt.Errorf("userRepo.GetActiveSegmentsIdByName: %w", err)
	}

	if len(segmentsId) != len(req.Segments) {
		return apperror.ErrNoSegment
	}

	err = s.userRepo.AddSegmentToUser(ctx, req.UserId, segmentsId, req.Ttl)
	if err != nil {
		return fmt.Errorf("userRepo.AddSegmentToUser: %w", err)
	}

	return nil
}

func (s *UserService) RemoveSegment(ctx context.Context, req entity.UserRemoveFromSegmentRequest) error {
	err := s.userRepo.CheckExistUser(ctx, req.UserId)
	if err != nil {
		return fmt.Errorf("userRepo.CheckExistUser: %w", err)
	}

	segmentsId, err := s.userRepo.GetActiveSegmentsIdByName(ctx, req.Segments)
	if err != nil {
		return fmt.Errorf("userRepo.GetActiveSegmentsIdByName: %w", err)
	}

	err = s.userRepo.RemoveSegmentFromUser(ctx, req.UserId, segmentsId)
	if err != nil {
		return fmt.Errorf("userRepo.RemoveSegmentFromUser: %w", err)
	}

	return nil
}

func (s *UserService) GetActiveSegments(ctx context.Context, req entity.UserActiveSegmentRequest) ([]string, error) {
	err := s.userRepo.CheckExistUser(ctx, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("userRepo.CheckExistUser: %w", err)
	}

	segments, err := s.userRepo.GetActiveSegmentFromUser(ctx, req.UserId)
	if err != nil {
		return nil, fmt.Errorf("userRepo.GetActiveSegmentFromUser: %w", err)
	}

	return segments, nil
}
