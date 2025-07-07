package service

import (
	"context"
	"term-service/internal/term/model"
	"term-service/internal/term/repository"
)

type TermService interface {
	CreateTerm(ctx context.Context, term *model.Term) (*model.Term, error)
	GetTermByID(ctx context.Context, id string) (*model.Term, error)
	UpdateTerm(ctx context.Context, id string, term *model.Term) error
	DeleteTerm(ctx context.Context, id string) error
	ListTerms(ctx context.Context) ([]*model.Term, error)
	GetCurrentTerm(ctx context.Context) (*model.Term, error)
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
