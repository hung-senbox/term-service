package mapper

import (
	"term-service/internal/holiday/dto/response"
	"term-service/internal/holiday/model"
	"term-service/pkg/helper"
)

func MapHolidayToResDTO(holiday *model.Holiday) response.HolidayResDTO {
	return response.HolidayResDTO{
		ID:               holiday.ID.Hex(),
		Color:            holiday.Color,
		PublishedMobile:  holiday.PublishedMobile,
		PublishedDesktop: holiday.PublishedDesktop,
		StartDate:        helper.FormatDate(holiday.StartDate),
		EndDate:          helper.FormatDate(holiday.EndDate),
		CreatedAt:        helper.FormatDate(holiday.CreatedAt),
	}
}

func MapHolidayListToResDTO(holidays []*model.Holiday) []response.HolidayResDTO {
	result := make([]response.HolidayResDTO, 0, len(holidays)) // slice rỗng, không phải nil
	for _, hld := range holidays {
		result = append(result, MapHolidayToResDTO(hld))
	}
	return result
}
