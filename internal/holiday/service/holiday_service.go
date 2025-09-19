package service

import (
	"context"
	"errors"
	"fmt"
	"term-service/internal/gateway"
	"term-service/internal/holiday/dto/request"
	"term-service/internal/holiday/dto/response"
	"term-service/internal/holiday/mapper"
	"term-service/internal/holiday/model"
	"term-service/internal/holiday/repository"
	pkg_helpder "term-service/pkg/helper"
	"time"

	"gorm.io/gorm"
)

type HolidayService interface {
	UploadHolidays(ctx context.Context, req request.UploadHolidayRequest) error
	GetHolidays4Web(ctx context.Context) (*response.GetHolidays4WebResDTO, error)
}

type holidayService struct {
	repo        repository.HolidayRepository
	userGateway gateway.UserGateway
	orgGateway  gateway.OrganizationGateway
}

func NewHolidayService(repo repository.HolidayRepository, userGateway gateway.UserGateway, orgGateway gateway.OrganizationGateway) HolidayService {
	return &holidayService{
		repo:        repo,
		userGateway: userGateway,
		orgGateway:  orgGateway,
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
		if err := s.repo.Delete(ctx, id); err != nil {
			return fmt.Errorf("failed to delete holiday %s: %w", id, err)
		}
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
		} else {
			// Create new Holiday
			newHoliday := &model.Holiday{
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

	if currentUser.IsSuperAdmin {
		// Lấy toàn bộ org từ Gateway
		orgs, err := s.orgGateway.GetAllOrg(ctx)
		if err != nil {
			return nil, fmt.Errorf("get all organizations failed: %w", err)
		}

		for _, org := range orgs {
			holidays, err := s.repo.GetAllByOrgID(ctx, org.ID)
			if err != nil {
				return nil, fmt.Errorf("get holidays by orgID %s failed: %w", org.ID, err)
			}

			result = append(result, response.HolidaysByOrgRes{
				OrganizationName: org.OrganizationName,
				Holidays:         mapper.MapHolidayListToResDTO(holidays),
			})
		}

	} else if currentUser.OrganizationAdmin.ID != "" {
		// User là org admin → chỉ lấy org của mình
		orgID := currentUser.OrganizationAdmin.ID
		holidays, err := s.repo.GetAllByOrgID(ctx, orgID)
		if err != nil {
			return nil, fmt.Errorf("get holidays by orgID %s failed: %w", orgID, err)
		}

		orgInfo, err := s.orgGateway.GetOrganizationInfo(ctx, orgID)
		if err != nil {
			return nil, fmt.Errorf("get organization info failed: %w", err)
		}

		result = append(result, response.HolidaysByOrgRes{
			OrganizationName: orgInfo.OrganizationName,
			Holidays:         mapper.MapHolidayListToResDTO(holidays),
		})

	} else {
		return nil, fmt.Errorf("access denied: user is not an organization admin")
	}

	return &response.GetHolidays4WebResDTO{
		HolidaysOrg: result,
	}, nil
}
