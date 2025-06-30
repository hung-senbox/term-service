package mappers

import (
	"term-info-service/internal/term/dto/response"
	"term-info-service/internal/term/model"
	"term-info-service/pkg/helper"
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
