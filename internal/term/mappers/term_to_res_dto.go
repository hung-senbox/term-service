package mappers

import (
	"term-service/internal/term/dto/response"
	"term-service/internal/term/model"
	"term-service/pkg/helper"
	"time"
)

func MapTermToResDTO(term *model.Term) response.TermResDTO {
	return response.TermResDTO{
		ID:        term.ID.Hex(),
		Title:     term.Title,
		StartDate: helper.FormatDate(term.StartDate),
		EndDate:   helper.FormatDate(term.EndDate),
		CreatedAt: helper.FormatDate(term.CreatedAt),
	}
}

func MapTermListToResDTO(terms []*model.Term) []response.TermResDTO {
	var result []response.TermResDTO
	for _, term := range terms {
		result = append(result, MapTermToResDTO(term))
	}
	return result
}

func MapTermToCurrentResDTO(term *model.Term) response.CurrentTermResDTO {
	layout := "2006-01-02"
	now := time.Now()
	remaining := int(term.EndDate.Sub(now).Hours() / 24)
	if remaining < 0 {
		remaining = 0
	}

	return response.CurrentTermResDTO{
		ID:           term.ID.Hex(),
		Title:        term.Title,
		StartDate:    term.StartDate.Format(layout),
		EndDate:      term.EndDate.Format(layout),
		CreatedAt:    term.CreatedAt.Format(layout),
		RemaningDate: helper.FormatRemainingDays(remaining),
	}
}
