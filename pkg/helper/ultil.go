package helper

import (
	"context"
	"strconv"
	"strings"
	gw_request "term-service/internal/gateway/dto/request"
	gw_response "term-service/internal/gateway/dto/response"
	"term-service/internal/holiday/dto/request"
	term_request "term-service/internal/term/dto/request"
	"term-service/pkg/constants"
)

func CurrentUserFromCtx(ctx context.Context) (*gw_response.CurrentUser, bool) {
	if cu, ok := ctx.Value(constants.CurrentUserKey).(*gw_response.CurrentUser); ok {
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

func BuildHolidayMessagesUpload(holidayID string, req request.UploadHolidayItem, langID uint) gw_request.UploadMessageLanguagesRequest {
	return gw_request.UploadMessageLanguagesRequest{
		MessageLanguages: []gw_request.UploadMessageRequest{
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

func BuildTermMessagesUpload(termID string, req term_request.UploadTermRequest, langID uint) gw_request.UploadMessageLanguagesRequest {
	return gw_request.UploadMessageLanguagesRequest{
		MessageLanguages: []gw_request.UploadMessageRequest{
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

func GetHeaders(ctx context.Context) map[string]string {
	headers := make(map[string]string)

	if lang, ok := ctx.Value(constants.AppLanguage).(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	return headers
}

func GetAppLanguage(ctx context.Context, defaultVal uint) uint {
	if lang, ok := ctx.Value(constants.AppLanguage).(uint); ok {
		return lang
	}
	return defaultVal
}
