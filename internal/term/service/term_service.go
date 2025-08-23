package service

import (
	"context"
	"errors"
	"fmt"
	"term-service/internal/gateway"
	"term-service/internal/term/dto/request"
	"term-service/internal/term/dto/response"
	"term-service/internal/term/mappers"
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
	ListTerms(ctx context.Context) (*response.ListTermsOrgResDTO, error)
	GetCurrentTerm(ctx context.Context) (*model.Term, error)
	UploadTerms(ctx context.Context, req []request.UploadTermItem) error
	GetTermsByOrgID(ctx context.Context, orgID string) (*response.ListTermsResDTO, error)
}

type termService struct {
	repo        repository.TermRepository
	userGateway gateway.UserGateway
	orgGateway  gateway.OrganizationGateway
}

func NewTermService(repo repository.TermRepository, userGateway gateway.UserGateway, orgGateway gateway.OrganizationGateway) TermService {
	return &termService{
		repo:        repo,
		userGateway: userGateway,
		orgGateway:  orgGateway,
	}
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

func (s *termService) ListTerms(ctx context.Context) (*response.ListTermsOrgResDTO, error) {
	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current user info failed: %w", err)
	}

	var terms []*model.Term
	var orgIDs []string

	if currentUser.IsSuperAdmin {
		terms, err = s.repo.GetAll(ctx)
		if err != nil {
			return nil, err
		}
		// gom tất cả orgID từ terms
		orgIDMap := map[string]struct{}{}
		for _, t := range terms {
			orgIDMap[t.OrganizationID] = struct{}{}
		}
		for id := range orgIDMap {
			orgIDs = append(orgIDs, id)
		}
	} else if currentUser.OrganizationAdmin.ID != "" {
		orgIDs = []string{currentUser.OrganizationAdmin.ID}
		terms, err = s.repo.GetAllByOrgID(ctx, currentUser.OrganizationAdmin.ID)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("access denied: user is not an organization admin")
	}

	// group terms by orgID
	termsByOrg := map[string][]*model.Term{}
	for _, t := range terms {
		termsByOrg[t.OrganizationID] = append(termsByOrg[t.OrganizationID], t)
	}

	// build response
	var result []response.TemsByOrgRes
	for _, orgID := range orgIDs {
		orgInfo, err := s.orgGateway.GetOrganizationInfo(ctx, orgID)
		if err != nil {
			return nil, fmt.Errorf("get organization info failed: %w", err)
		}

		result = append(result, response.TemsByOrgRes{
			OrganizationName: orgInfo.OrganizationName, // chỉ lấy tên
			Terms:            mappers.MapTermListToResDTO(termsByOrg[orgID]),
		})

	}

	return &response.ListTermsOrgResDTO{
		TermsOrg: result,
	}, nil
}

func (s *termService) GetCurrentTerm(ctx context.Context) (*model.Term, error) {
	return s.repo.GetCurrentTerm(ctx)
}

func (s *termService) UploadTerms(ctx context.Context, req []request.UploadTermItem) error {
	// get organization admin from user context
	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil {
		return fmt.Errorf("get current user info failed")
	}

	// check is super admin & check org admin
	if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
		logger.WriteLogEx("warn", "Access denied for super admin", map[string]any{
			"user_id": currentUser.ID,
		})
		return fmt.Errorf("access denied: super admin cannot perform this action")
	}
	organizationAdminID := currentUser.OrganizationAdmin.ID

	// Upsert terms (không cần delete nữa)
	for _, t := range req {
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
			existing.PublishedTeacher = t.PublishedTeacher
			existing.PublishedParent = t.PublishedParent
			existing.StartDate = startDate
			existing.EndDate = endDate
			existing.UpdatedAt = time.Now()

			if err := s.repo.Update(ctx, t.ID, existing); err != nil {
				return fmt.Errorf("failed to update term %s: %w", t.ID, err)
			}
		} else {
			// Create new term
			newTerm := &model.Term{
				OrganizationID:   organizationAdminID,
				Title:            t.Title,
				Color:            t.Color,
				PublishedMobile:  t.PublishedMobile,
				PublishedDesktop: t.PublishedDesktop,
				PublishedTeacher: t.PublishedTeacher,
				PublishedParent:  t.PublishedParent,
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

func (s *termService) GetTermsByOrgID(ctx context.Context, orgID string) (*response.ListTermsResDTO, error) {
	terms, err := s.repo.GetAllByOrgID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("get terms by orgID failed: %w", err)
	}

	if len(terms) == 0 {
		return &response.ListTermsResDTO{
			Terms: make([]response.TermResDTO, 0),
		}, nil
	}

	return &response.ListTermsResDTO{
		Terms: mappers.MapTermListToResDTO(terms),
	}, nil
}
