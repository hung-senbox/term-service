package service

import (
	"context"
	"errors"
	"fmt"
	"term-service/internal/gateway"
	"term-service/internal/gateway/dto"
	"term-service/internal/holiday/dto/request"
	"term-service/internal/holiday/dto/response"
	"term-service/internal/holiday/mapper"
	"term-service/internal/holiday/model"
	"term-service/internal/holiday/repository"
	"term-service/pkg/constants"
	"term-service/pkg/helper"
	pkg_helpder "term-service/pkg/helper"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

type HolidayService interface {
	UploadHolidays(ctx context.Context, req request.UploadHolidayRequest) error
	GetHolidays4Web(ctx context.Context) (*response.GetHolidays4WebResDTO, error)
}

type holidayService struct {
	repo                   repository.HolidayRepository
	userGateway            gateway.UserGateway
	orgGateway             gateway.OrganizationGateway
	messageLanguageGateway gateway.MessageLanguageGateway
}

func NewHolidayService(repo repository.HolidayRepository, userGateway gateway.UserGateway, orgGateway gateway.OrganizationGateway, messageLanguageGateway gateway.MessageLanguageGateway) HolidayService {
	return &holidayService{
		repo:                   repo,
		userGateway:            userGateway,
		orgGateway:             orgGateway,
		messageLanguageGateway: messageLanguageGateway,
	}
}

func (s *holidayService) UploadHolidays(ctx context.Context, req request.UploadHolidayRequest) error {
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

	// 1. Handle delete
	for _, id := range req.DeleteIds {
		err := s.repo.Delete(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to delete holiday")
		}
		// goi GW xoa message lang
		s.messageLanguageGateway.DeleleByTypeAndTypeID(ctx, string(constants.HolidayType), id)
	}

	// 2. Handle upsert (create or update)
	for _, t := range req.Holidays {
		startDate, err := time.Parse("2006-01-02", t.StartDate)
		if err != nil {
			return fmt.Errorf("invalid start_date for holiday %s: %w", t.Title, err)
		}

		endDate, err := time.Parse("2006-01-02", t.EndDate)
		if err != nil {
			return fmt.Errorf("invalid end_date for holiday %s: %w", t.Title, err)
		}

		if !pkg_helpder.ValidateDateRange(startDate, endDate) {
			return fmt.Errorf("start_date must be before or equal to end_date for holiday %s", t.Title)
		}

		if t.ID != "" {
			// Update existing holiday
			existing, err := s.repo.GetByID(ctx, t.ID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("holiday not found: %s", t.ID)
				}
				return fmt.Errorf("failed to get holiday: %w", err)
			}

			existing.Title = t.Title
			existing.Color = t.Color
			existing.PublishedMobile = t.PublishedMobile
			existing.PublishedDesktop = t.PublishedDesktop
			existing.StartDate = startDate
			existing.EndDate = endDate
			existing.UpdatedAt = time.Now()

			if err := s.repo.Update(ctx, t.ID, existing); err != nil {
				return fmt.Errorf("failed to update holiday %s: %w", t.ID, err)
			}

			// goi messs lang gw upload message
			err = s.uploadMessages(ctx, helper.BuildHolidayMessagesUpload(existing.ID.Hex(), t, req.LanguageID))
			if err != nil {
				return fmt.Errorf("upload department messages failed")
			}

		} else {
			// Create new Holiday
			newHoliday := &model.Holiday{
				ID:               primitive.NewObjectID(),
				OrganizationID:   organizationAdminID,
				Title:            t.Title,
				Color:            t.Color,
				PublishedMobile:  t.PublishedMobile,
				PublishedDesktop: t.PublishedDesktop,
				StartDate:        startDate,
				EndDate:          endDate,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}

			if _, err := s.repo.Create(ctx, newHoliday); err != nil {
				return fmt.Errorf("failed to create holiday %s: %w", t.Title, err)
			}

			// goi messs lang gw upload message
			err = s.uploadMessages(ctx, helper.BuildHolidayMessagesUpload(newHoliday.ID.Hex(), t, req.LanguageID))

			if err != nil {
				return fmt.Errorf("upload department messages failed")
			}
		}
	}

	return nil
}

func (s *holidayService) GetHolidays4Web(ctx context.Context) (*response.GetHolidays4WebResDTO, error) {
	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current user info failed: %w", err)
	}

	var result []response.HolidaysByOrgRes

	// if is super admin return []
	if currentUser.IsSuperAdmin {
		// // Lấy toàn bộ org từ Gateway
		// orgs, err := s.orgGateway.GetAllOrg(ctx)
		// if err != nil {
		// 	return nil, fmt.Errorf("get all organizations failed: %w", err)
		// }

		// for _, org := range orgs {
		// 	holidays, err := s.repo.GetAllByOrgID(ctx, org.ID)
		// 	if err != nil {
		// 		return nil, fmt.Errorf("get holidays by orgID %s failed: %w", org.ID, err)
		// 	}
		// 	holidayDTOs := mapper.MapHolidayListToResDTO(holidays)

		// 	// --- bổ sung message languages ---
		// 	for i := range holidayDTOs {
		// 		msgLangs, _ := s.messageLanguageGateway.GetMessageLanguages(ctx, "holiday", holidayDTOs[i].ID)
		// 		if msgLangs == nil {
		// 			msgLangs = []dto.MessageLanguageResponse{}
		// 		}
		// 		holidayDTOs[i].MessageLanguages = msgLangs
		// 	}

		// 	result = append(result, response.HolidaysByOrgRes{
		// 		OrganizationName: org.OrganizationName,
		// 		Holidays:         holidayDTOs,
		// 	})
		// }

		return &response.GetHolidays4WebResDTO{
			HolidaysOrg: result,
		}, nil

	} else if currentUser.OrganizationAdmin.ID != "" {
		// User là org admin → chỉ lấy org của mình
		orgID := currentUser.OrganizationAdmin.ID
		holidays, err := s.repo.GetAllByOrgID(ctx, orgID)
		if err != nil {
			return nil, fmt.Errorf("get holidays by orgID %s failed: %w", orgID, err)
		}

		holidayDTOs := mapper.MapHolidayListToResDTO(holidays)

		// --- bổ sung message languages ---
		for i := range holidayDTOs {
			msgLangs, _ := s.messageLanguageGateway.GetMessageLanguages(ctx, "holiday", holidayDTOs[i].ID)
			if msgLangs == nil {
				msgLangs = []dto.MessageLanguageResponse{}
			}
			holidayDTOs[i].MessageLanguages = msgLangs
		}

		orgInfo, err := s.orgGateway.GetOrganizationInfo(ctx, orgID)
		if err != nil {
			return nil, fmt.Errorf("get organization info failed: %w", err)
		}

		result = append(result, response.HolidaysByOrgRes{
			OrganizationName: orgInfo.OrganizationName,
			Holidays:         holidayDTOs,
		})

	} else {
		return nil, fmt.Errorf("access denied: user is not an organization admin")
	}

	return &response.GetHolidays4WebResDTO{
		HolidaysOrg: result,
	}, nil
}

func (s *holidayService) uploadMessages(ctx context.Context, req dto.UploadMessageLanguagesRequest) error {
	err := s.messageLanguageGateway.UploadMessages(ctx, req)

	if err != nil {
		return err
	}

	return nil
}
