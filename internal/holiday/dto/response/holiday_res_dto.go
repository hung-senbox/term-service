package response

import "term-service/internal/gateway/dto/response"

type HolidayResDTO struct {
	ID               string                             `json:"id"`
	Color            string                             `json:"color"`
	PublishedMobile  bool                               `json:"published_mobile"`
	PublishedDesktop bool                               `json:"published_desktop"`
	StartDate        string                             `json:"start_date"`
	EndDate          string                             `json:"end_date"`
	CreatedAt        string                             `json:"created_at"`
	MessageLanguages []response.MessageLanguageResponse `json:"message_languages"`
}
