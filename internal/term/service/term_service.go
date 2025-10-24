package service

import (
	"context"
	"errors"
	"fmt"
	"term-service/internal/gateway"
	gw_request "term-service/internal/gateway/dto/request"
	gw_response "term-service/internal/gateway/dto/response"
	"term-service/internal/term/dto/request"
	"term-service/internal/term/dto/response"
	"term-service/internal/term/mappers"
	"term-service/internal/term/model"
	"term-service/internal/term/repository"
	pkg_helpder "term-service/pkg/helper"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

type TermService interface {
	CreateTerm(ctx context.Context, term *model.Term) (*model.Term, error)
	GetTermByID(ctx context.Context, id string) (*model.Term, error)
	UpdateTerm(ctx context.Context, id string, term *model.Term) error
	DeleteTerm(ctx context.Context, id string) error
	GetTerms4Web(ctx context.Context) (*response.GetTerms4WebResDTO, error)
	GetCurrentTerm(ctx context.Context) (response.CurrentTermResDTO, error)
	UploadTerms(ctx context.Context, req request.UploadTermRequest) error
	GetTermsByOrgID(ctx context.Context, orgID string) (*response.ListTermsResDTO, error)
	GetTermsByStudent4App(ctx context.Context, studentID string) ([]response.TermsByStudentResDTO, error)
	GetTermsByStudent4Web(ctx context.Context, studentID string) ([]response.TermsByStudentResDTO, error)
	GetCurrentTermByOrg(ctx context.Context, organizationID string) (response.CurrentTermResDTO, error)
	GetTerms4App(ctx context.Context, organizationID string) (*response.GetTerms4AppResDTO, error)
	GetTerm4Gw(ctx context.Context, termId string) (*response.Term4GwResponse, error)
	GetTermsByOrg4App(ctx context.Context, organizationID string) ([]response.TermResponse4App, error)
	GetPreviousTerm4GW(ctx context.Context, organizationID string, termID string) (*response.Term4GwResponse, error)
	GetPreviousTerms4GW(ctx context.Context, organizationID string, termID string) ([]*response.Term4GwResponse, error)
	GetTerms2Assign4Web(ctx context.Context) ([]*response.TermResponse4Web, error)
}

type termService struct {
	repo                   repository.TermRepository
	cachedUserGw           gateway.UserGateway
	orgGateway             gateway.OrganizationGateway
	messageLanguageGateway gateway.MessageLanguageGateway
}

func NewTermService(
	repo repository.TermRepository,
	cachedUserGw gateway.UserGateway,
	orgGateway gateway.OrganizationGateway,
	messageLanguageGateway gateway.MessageLanguageGateway,
) TermService {
	return &termService{
		repo:                   repo,
		cachedUserGw:           cachedUserGw,
		orgGateway:             orgGateway,
		messageLanguageGateway: messageLanguageGateway,
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
	currentUser, err := s.cachedUserGw.GetCurrentUser(ctx)
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

			// --- gọi message language gateway ---
			msgLangs, _ := s.messageLanguageGateway.GetMessageLanguages(ctx, "term", org.ID)
			if msgLangs == nil {
				msgLangs = []gw_response.MessageLanguageResponse{}
			}

			result = append(result, response.TermsByOrgRes{
				MessageLanguages: msgLangs,
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

		// --- gọi message language gateway ---
		msgLangs, _ := s.messageLanguageGateway.GetMessageLanguages(ctx, "term", orgID)
		if msgLangs == nil {
			msgLangs = []gw_response.MessageLanguageResponse{}
		}

		result = append(result, response.TermsByOrgRes{
			MessageLanguages: msgLangs,
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

	return mappers.MapTermToCurrentResDTO(term, ""), nil
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

	return mappers.MapTermToCurrentResDTO(term, ""), nil
}

func (s *termService) UploadTerms(ctx context.Context, req request.UploadTermRequest) error {
	// get organization admin from user context
	currentUser, err := s.cachedUserGw.GetCurrentUser(ctx)
	if err != nil {
		return fmt.Errorf("get current user info failed")
	}

	// check is super admin & check org admin
	if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
		return fmt.Errorf("access denied: super admin cannot perform this action")
	}
	organizationAdminID := currentUser.OrganizationAdmin.ID

	// Upsert terms
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

			// goi messs lang gw upload message
			err = s.uploadMessages(ctx, pkg_helpder.BuildTermMessagesUpload(organizationAdminID, req, req.LanguageID))
			if err != nil {
				return fmt.Errorf("upload department messages failed")
			}

		} else {
			// Create new term
			newTerm := &model.Term{
				ID:               primitive.NewObjectID(),
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

			err = s.uploadMessages(ctx, pkg_helpder.BuildTermMessagesUpload(organizationAdminID, req, req.LanguageID))
			if err != nil {
				return fmt.Errorf("upload department messages failed")
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

func (s *termService) GetTermsByStudent4App(ctx context.Context, studentID string) ([]response.TermsByStudentResDTO, error) {
	// get student info
	student, err := s.cachedUserGw.GetStudentInfo(ctx, studentID)
	if err != nil {
		return nil, fmt.Errorf("get student info failed: %w", err)
	}

	// get terms by orgID
	terms, err := s.repo.GetAllByOrgID4App(ctx, student.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("get terms by orgID failed: %w", err)
	}

	if len(terms) == 0 {
		return []response.TermsByStudentResDTO{}, nil
	}

	// get word by orgID
	msg, _ := s.messageLanguageGateway.GetMessageLanguage(ctx, "term", student.OrganizationID)
	word := ""
	if msg.Contents != nil {
		if val, ok := msg.Contents["word"]; ok {
			word = val
		}
	}

	return mappers.MapTermsByStudentToResDTO(terms, word), nil
}

func (s *termService) GetTermsByStudent4Web(ctx context.Context, studentID string) ([]response.TermsByStudentResDTO, error) {
	// get student info
	student, err := s.cachedUserGw.GetStudentInfo(ctx, studentID)
	if err != nil {
		return nil, fmt.Errorf("get student info failed: %w", err)
	}

	// get terms by orgID
	terms, err := s.repo.GetAllByOrgID4Web(ctx, student.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("get terms by orgID failed: %w", err)
	}

	if len(terms) == 0 {
		return []response.TermsByStudentResDTO{}, nil
	}

	// get word by orgID
	msg, _ := s.messageLanguageGateway.GetMessageLanguage(ctx, "term", student.OrganizationID)
	word := ""
	if msg.Contents != nil {
		if val, ok := msg.Contents["word"]; ok {
			word = val
		}
	}

	return mappers.MapTermsByStudentToResDTO(terms, word), nil
}

func (s *termService) GetTerms4App(ctx context.Context, organizationID string) (*response.GetTerms4AppResDTO, error) {
	terms, err := s.repo.GetAllByOrgID4App(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("get terms by orgID failed: %w", err)
	}

	// get word by orgID
	msg, _ := s.messageLanguageGateway.GetMessageLanguage(ctx, "term", organizationID)
	word := ""
	if msg.Contents != nil {
		if val, ok := msg.Contents["word"]; ok {
			word = val
		}
	}

	return &response.GetTerms4AppResDTO{
		Terms: mappers.MapTermListToCurrentResDTO(terms, word),
	}, nil
}

func (s *termService) uploadMessages(ctx context.Context, req gw_request.UploadMessageLanguagesRequest) error {
	err := s.messageLanguageGateway.UploadMessages(ctx, req)

	if err != nil {
		return err
	}

	return nil
}

func (s *termService) GetTerm4Gw(ctx context.Context, termId string) (*response.Term4GwResponse, error) {
	// get organization admin from user context
	currentUser, err := s.cachedUserGw.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current user info failed")
	}

	term, err := s.repo.GetByID(ctx, termId)
	if err != nil {
		return nil, fmt.Errorf("get term by id failed: %w", err)
	}

	// get word by orgID
	msg, _ := s.messageLanguageGateway.GetMessageLanguage(ctx, "term", currentUser.OrganizationAdmin.ID)
	word := ""
	if msg.Contents != nil {
		if val, ok := msg.Contents["word"]; ok {
			word = val
		}
	}

	return mappers.MapTermToRes4GwResponse(term, word), nil
}

func (s *termService) GetTermsByOrg4App(ctx context.Context, organizationID string) ([]response.TermResponse4App, error) {

	// get terms by orgID
	terms, err := s.repo.GetAllByOrgIDIsPublishedTeacher(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("get terms by orgID failed: %w", err)
	}

	if len(terms) == 0 {
		return []response.TermResponse4App{}, nil
	}

	// get word by orgID
	msg, _ := s.messageLanguageGateway.GetMessageLanguage(ctx, "term", organizationID)
	word := ""
	if msg.Contents != nil {
		if val, ok := msg.Contents["word"]; ok {
			word = val
		}
	}

	return mappers.MapTermsByToRes4App(terms, word), nil
}

func (s *termService) GetPreviousTerm4GW(ctx context.Context, organizationID string, termID string) (*response.Term4GwResponse, error) {

	// get terms by orgID
	previousTerm, err := s.repo.GetPreviousTerm(ctx, organizationID, termID)
	if err != nil {
		return nil, fmt.Errorf("get terms by orgID failed: %w", err)
	}

	if previousTerm == nil {
		return nil, fmt.Errorf("get terms by orgID failed: %w", err)
	}

	// get word by orgID
	msg, _ := s.messageLanguageGateway.GetMessageLanguage(ctx, "term", organizationID)
	word := ""
	if msg.Contents != nil {
		if val, ok := msg.Contents["word"]; ok {
			word = val
		}
	}

	return mappers.MapTermToRes4GwResponse(previousTerm, word), nil
}

func (s *termService) GetPreviousTerms4GW(ctx context.Context, organizationID string, termID string) ([]*response.Term4GwResponse, error) {

	// get terms by orgID
	previousTerms, err := s.repo.GetPreviousTerms(ctx, organizationID, termID)
	if err != nil {
		return nil, fmt.Errorf("get terms by orgID failed: %w", err)
	}

	// get word by orgID
	msg, _ := s.messageLanguageGateway.GetMessageLanguage(ctx, "term", organizationID)
	word := ""
	if msg.Contents != nil {
		if val, ok := msg.Contents["word"]; ok {
			word = val
		}
	}

	return mappers.MapTermsToRes4GwResponse(previousTerms, word), nil
}

func (s *termService) GetTerms2Assign4Web(ctx context.Context) ([]*response.TermResponse4Web, error) {
	// get organization admin from user context
	currentUser, err := s.cachedUserGw.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current user info failed")
	}

	// check is super admin & check org admin
	if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
		return nil, fmt.Errorf("access denied: super admin cannot perform this action")
	}

	// get terms by orgID
	terms, err := s.repo.GetAllByOrgIDIsPublishedDesktop(ctx, currentUser.OrganizationAdmin.ID)
	if err != nil {
		return nil, fmt.Errorf("get terms by orgID failed: %w", err)
	}

	// get word by orgID
	msg, _ := s.messageLanguageGateway.GetMessageLanguage(ctx, "term", currentUser.OrganizationAdmin.ID)
	word := ""
	if msg.Contents != nil {
		if val, ok := msg.Contents["word"]; ok {
			word = val
		}
	}

	var res = make([]*response.TermResponse4Web, 0, len(terms))

	for _, term := range terms {
		res = append(res, mappers.MapTermToRes4WebResponse(term, word))
	}

	return res, nil
}
