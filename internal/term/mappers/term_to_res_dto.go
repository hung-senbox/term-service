package mappers

import (
	"term-service/internal/term/dto/response"
	"term-service/internal/term/model"
	"term-service/pkg/helper"
	"time"
)

func MapTermToResDTO(term *model.Term) response.TermResDTO {
	return response.TermResDTO{
		ID:               term.ID.Hex(),
		Title:            term.Title,
		Color:            term.Color,
		PublishedMobile:  term.PublishedMobile,
		PublishedDesktop: term.PublishedDesktop,
		StartDate:        helper.FormatDate(term.StartDate),
		EndDate:          helper.FormatDate(term.EndDate),
		CreatedAt:        helper.FormatDate(term.CreatedAt),
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
	now := time.Now().In(term.EndDate.Location())

	// get remning days
	remaining := daysBetweenDateOnly(now, term.EndDate)
	if remaining < 0 {
		remaining = 0
	}
	// gert current wweek
	currentWeek := calculateCurrentWeek(term.StartDate, term.EndDate, now)

	return response.CurrentTermResDTO{
		ID:           term.ID.Hex(),
		Title:        term.Title,
		Color:        term.Color,
		StartDate:    term.StartDate.Format(layout),
		EndDate:      term.EndDate.Format(layout),
		CreatedAt:    term.CreatedAt.Format(layout),
		RemaningDate: helper.FormatRemainingDays(remaining),
		CurrentWeek:  currentWeek,
	}
}

func daysBetweenDateOnly(start, end time.Time) int {
	loc := end.Location()
	sy, sm, sd := start.In(loc).Date()
	ey, em, ed := end.In(loc).Date()

	startDate := time.Date(sy, sm, sd, 0, 0, 0, 0, loc)
	endDate := time.Date(ey, em, ed, 0, 0, 0, 0, loc)

	return int(endDate.Sub(startDate).Hours() / 24)
}

func calculateCurrentWeek(start, end, now time.Time) int {
	loc := end.Location()
	sy, sm, sd := start.In(loc).Date()
	ey, em, ed := end.In(loc).Date()
	ny, nm, nd := now.In(loc).Date()

	startDate := time.Date(sy, sm, sd, 0, 0, 0, 0, loc)
	endDate := time.Date(ey, em, ed, 0, 0, 0, 0, loc)
	nowDate := time.Date(ny, nm, nd, 0, 0, 0, 0, loc)

	if nowDate.Before(startDate) {
		return 0
	}
	if nowDate.After(endDate) {
		totalDays := int(endDate.Sub(startDate).Hours() / 24)
		return (totalDays / 7) + 1
	}

	daysPassed := int(nowDate.Sub(startDate).Hours() / 24)
	return (daysPassed / 7) + 1
}
