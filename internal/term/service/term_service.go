package service

import (
	"context"
	"errors"
	"fmt"
	"term-service/internal/term/dto/request"
	"term-service/internal/term/model"
	"term-service/internal/term/repository"
	"term-service/logger"
	pkg_helpder "term-service/pkg/helper"
	"time"

	"gorm.io/gorm"
)

type TermService interface {
	CreateTerm(ctx context.Context, term *model.Term) (*model.Term, error)
	GetTermByID(ctx context.Context, id string) (*model.Term, error)
	UpdateTerm(ctx context.Context, id string, term *model.Term) error
	DeleteTerm(ctx context.Context, id string) error
	ListTerms(ctx context.Context) ([]*model.Term, error)
	GetCurrentTerm(ctx context.Context) (*model.Term, error)
	UploadTerms(ctx context.Context, req *request.UploadTermReqDTO) error
}

type termService struct {
	repo repository.TermRepository
}

func NewTermService(repo repository.TermRepository) TermService {
	return &termService{repo: repo}
}

func (s *termService) CreateTerm(ctx context.Context, term *model.Term) (*model.Term, error) {
	return s.repo.Create(ctx, term)
}

func (s *termService) GetTermByID(ctx context.Context, id string) (*model.Term, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *termService) UpdateTerm(ctx context.Context, id string, term *model.Term) error {
	return s.repo.Update(ctx, id, term)
}

func (s *termService) DeleteTerm(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *termService) ListTerms(ctx context.Context) ([]*model.Term, error) {
	return s.repo.GetAll(ctx)
}

func (s *termService) GetCurrentTerm(ctx context.Context) (*model.Term, error) {
	return s.repo.GetCurrentTerm(ctx)
}

func (s *termService) UploadTerms(ctx context.Context, req *request.UploadTermReqDTO) error {
	// 1. Delete terms if any IDs provided
	if len(req.DeleteTermIds) > 0 {
		for _, id := range req.DeleteTermIds {
			if err := s.repo.Delete(ctx, id); err != nil {
				return fmt.Errorf("failed to delete term %s: %w", id, err)
			}
		}
	}

	// 2. Upsert terms
	for _, t := range req.Terms {
		startDate, err := time.Parse("2006-01-02", t.StartDate)
		if err != nil {
			return fmt.Errorf("invalid start_date for term %s: %w", t.Title, err)
		}

		endDate, err := time.Parse("2006-01-02", t.EndDate)
		if err != nil {
			return fmt.Errorf("invalid end_date for term %s: %w", t.Title, err)
		}

		if !pkg_helpder.ValidateDateRange(startDate, endDate) {
			return fmt.Errorf("start_date must be before or equal to end_date for term %s", t.Title)
		}

		if t.ID != "" {
			// Update existing term
			existing, err := s.repo.GetByID(ctx, t.ID)
			if err != nil {
				// Log chi tiết ở service
				logger.WriteLogEx("error", "Get term in UploadTerms failed", map[string]any{
					"term_id": t.ID,
					"error":   err.Error(),
				})
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("term not found: %s", t.ID)
				}
				return fmt.Errorf("failed to get term: %w", err)
			}
			existing.Title = t.Title
			existing.Color = t.Color
			existing.PublishedMobile = t.PublishedMobile
			existing.PublishedDesktop = t.PublishedDesktop
			existing.StartDate = startDate
			existing.EndDate = endDate
			existing.UpdatedAt = time.Now()

			if err := s.repo.Update(ctx, t.ID, existing); err != nil {
				return fmt.Errorf("failed to update term %s: %w", t.ID, err)
			}
		} else {
			// Create new term
			newTerm := &model.Term{
				Title:            t.Title,
				Color:            t.Color,
				PublishedMobile:  t.PublishedMobile,
				PublishedDesktop: t.PublishedDesktop,
				StartDate:        startDate,
				EndDate:          endDate,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}
			if _, err := s.repo.Create(ctx, newTerm); err != nil {
				return fmt.Errorf("failed to create term %s: %w", t.Title, err)
			}
		}
	}

	return nil
}
