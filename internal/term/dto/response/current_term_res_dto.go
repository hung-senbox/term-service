package response

type CurrentTermResDTO struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	CreatedAt    string `json:"created_at"`
	RemaningDate string `json:"remaning_date"`
}
