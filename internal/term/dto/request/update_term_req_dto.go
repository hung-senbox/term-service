package request

type UpdateTermReqDTO struct {
	Title     *string `json:"title"`      // optional
	StartDate *string `json:"start_date"` // optional, format: YYYY-MM-DD
	EndDate   *string `json:"end_date"`   // optional, format: YYYY-MM-DD
}
