package helper

import (
	"context"
	"strconv"
	"strings"
	"term-service/internal/gateway/dto"
	"term-service/internal/holiday/dto/request"
	term_request "term-service/internal/term/dto/request"
	"term-service/pkg/constants"
)

func CurrentUserFromCtx(ctx context.Context) (*dto.CurrentUser, bool) {
	if cu, ok := ctx.Value(constants.CurrentUserKey).(*dto.CurrentUser); ok {
		return cu, true
	}
	return nil, false
}

func ParseAppLanguage(header string, defaultVal uint) uint {
	header = strings.TrimSpace(strings.Trim(header, "\""))
	if val, err := strconv.Atoi(header); err == nil {
		return uint(val)
	}
	return defaultVal
}

func BuildHolidayMessagesUpload(holidayID string, req request.UploadHolidayItem, langID uint) dto.UploadMessageLanguagesRequest {
	return dto.UploadMessageLanguagesRequest{
		MessageLanguages: []dto.UploadMessageRequest{
			{
				TypeID:     holidayID,
				Type:       string(constants.HolidayType),
				Key:        string(constants.HolidayTitleKey),
				Value:      req.Title,
				LanguageID: langID,
			},
		},
	}
}

func BuildTermMessagesUpload(termID string, req term_request.UploadTermRequest, langID uint) dto.UploadMessageLanguagesRequest {
	return dto.UploadMessageLanguagesRequest{
		MessageLanguages: []dto.UploadMessageRequest{
			{
				TypeID:     termID,
				Type:       string(constants.TermType),
				Key:        string(constants.TermWordKey),
				Value:      req.Word,
				LanguageID: langID,
			},
		},
	}
}
