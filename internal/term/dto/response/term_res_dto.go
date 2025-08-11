package response

type TermResDTO struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Color     string `json:"color"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	CreatedAt string `json:"created_at"`
}
