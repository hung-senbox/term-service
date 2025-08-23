package request

type UploadTermItem struct {
	ID               string `json:"id,omitempty"`
	Title            string `json:"title"`
	Color            string `json:"color"`
	PublishedMobile  bool   `json:"published_mobile"`
	PublishedDesktop bool   `json:"published_desktop"`
	PublishedTeacher bool   `json:"published_teacher"`
	PublishedParent  bool   `json:"published_parent"`
	StartDate        string `json:"start_date"`
	EndDate          string `json:"end_date"`
}
