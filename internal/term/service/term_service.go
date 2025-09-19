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
	pkg_helpder "term-service/pkg/helper"
	"time"

	"gorm.io/gorm"
)

type TermService interface {
	CreateTerm(ctx context.Context, term *model.Term) (*model.Term, error)
	GetTermByID(ctx context.Context, id string) (*model.Term, error)
	UpdateTerm(ctx context.Context, id string, term *model.Term) error
	DeleteTerm(ctx context.Context, id string) error
	GetTerms4Web(ctx context.Context) (*response.GetTerms4WebResDTO, error)
	GetCurrentTerm(ctx context.Context) (response.CurrentTermResDTO, error)
	UploadTerms(ctx context.Context, req []request.UploadTermItem) error
	GetTermsByOrgID(ctx context.Context, orgID string) (*response.ListTermsResDTO, error)
	GetTermsByStudent(ctx context.Context, studentID string) ([]response.TermsByStudentResDTO, error)
	GetCurrentTermByOrg(ctx context.Context, organizationID string) (response.CurrentTermResDTO, error)
	GetTerms4App(ctx context.Context, organizationID string) (*response.GetTerms4AppResDTO, error)
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

func (s *termService) GetTerms4Web(ctx context.Context) (*response.GetTerms4WebResDTO, error) {
	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current user info failed: %w", err)
	}

	var result []response.TermsByOrgRes

	if currentUser.IsSuperAdmin {
		// Lấy toàn bộ org từ Gateway
		orgs, err := s.orgGateway.GetAllOrg(ctx)
		if err != nil {
			return nil, fmt.Errorf("get all organizations failed: %w", err)
		}

		for _, org := range orgs {
			terms, err := s.repo.GetAllByOrgID(ctx, org.ID)
			if err != nil {
				return nil, fmt.Errorf("get terms by orgID %s failed: %w", org.ID, err)
			}

			result = append(result, response.TermsByOrgRes{
				OrganizationName: org.OrganizationName,
				Terms:            mappers.MapTermListToResDTO(terms),
			})
		}

	} else if currentUser.OrganizationAdmin.ID != "" {
		// User là org admin → chỉ lấy org của mình
		orgID := currentUser.OrganizationAdmin.ID
		terms, err := s.repo.GetAllByOrgID(ctx, orgID)
		if err != nil {
			return nil, fmt.Errorf("get terms by orgID %s failed: %w", orgID, err)
		}

		orgInfo, err := s.orgGateway.GetOrganizationInfo(ctx, orgID)
		if err != nil {
			return nil, fmt.Errorf("get organization info failed: %w", err)
		}

		result = append(result, response.TermsByOrgRes{
			OrganizationName: orgInfo.OrganizationName,
			Terms:            mappers.MapTermListToResDTO(terms),
		})

	} else {
		return nil, fmt.Errorf("access denied: user is not an organization admin")
	}

	return &response.GetTerms4WebResDTO{
		TermsOrg: result,
	}, nil
}

func (s *termService) GetCurrentTerm(ctx context.Context) (response.CurrentTermResDTO, error) {
	term, err := s.repo.GetCurrentTerm(ctx)
	if err != nil {
		return response.CurrentTermResDTO{}, fmt.Errorf("get current term failed: %w", err)
	}

	if term == nil {
		return response.CurrentTermResDTO{}, fmt.Errorf("no current term found")
	}

	return mappers.MapTermToCurrentResDTO(term), nil
}

func (s *termService) GetCurrentTermByOrg(ctx context.Context, organizationID string) (response.CurrentTermResDTO, error) {
	var term *model.Term
	var err error

	if organizationID == "" {
		term, err = s.repo.GetCurrentTerm(ctx)
	} else {
		term, err = s.repo.GetCurrentTermByOrg(ctx, organizationID)
	}

	if err != nil {
		return response.CurrentTermResDTO{}, fmt.Errorf("get current term failed: %w", err)
	}

	if term == nil {
		return response.CurrentTermResDTO{}, fmt.Errorf("no current term found")
	}

	return mappers.MapTermToCurrentResDTO(term), nil
}

func (s *termService) UploadTerms(ctx context.Context, req []request.UploadTermItem) error {
	// get organization admin from user context
	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil {
		return fmt.Errorf("get current user info failed")
	}

	// check is super admin & check org admin
	if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
		return fmt.Errorf("access denied: super admin cannot perform this action")
	}
	organizationAdminID := currentUser.OrganizationAdmin.ID

	// Upsert
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

func (s *termService) GetTermsByStudent(ctx context.Context, studentID string) ([]response.TermsByStudentResDTO, error) {
	// get student info
	student, err := s.userGateway.GetStudentInfo(ctx, studentID)
	if err != nil {
		return nil, fmt.Errorf("get student info failed: %w", err)
	}

	// get terms by orgID
	terms, err := s.repo.GetAllByOrgID(ctx, student.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("get terms by orgID failed: %w", err)
	}

	if len(terms) == 0 {
		return []response.TermsByStudentResDTO{}, nil
	}

	return mappers.MapTermsByStudentToResDTO(terms), nil
}

func (s *termService) GetTerms4App(ctx context.Context, organizationID string) (*response.GetTerms4AppResDTO, error) {
	terms, err := s.repo.GetAllByOrgID4App(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("get terms by orgID failed: %w", err)
	}

	return &response.GetTerms4AppResDTO{
		Terms: mappers.MapTermListToCurrentResDTO(terms),
	}, nil
}
