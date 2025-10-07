package response

type TermResDTO struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Color            string `json:"color"`
	PublishedMobile  bool   `json:"published_mobile"`
	PublishedDesktop bool   `json:"published_desktop"`
	PublishedTeacher bool   `json:"published_teacher"`
	PublishedParent  bool   `json:"published_parent"`
	StartDate        string `json:"start_date"`
	EndDate          string `json:"end_date"`
	CreatedAt        string `json:"created_at"`
}

type TermsByStudentResDTO struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Color string `json:"color"`
}

type TermResponse4App struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Color string `json:"color"`
}
