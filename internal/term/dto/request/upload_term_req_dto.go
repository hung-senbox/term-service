package request

type UploadTermReqDTO struct {
	DeleteTermIds []string         `json:"delete_term_ids"`
	Terms         []UploadTermItem `json:"terms"`
}

type UploadTermItem struct {
	ID               string `json:"id,omitempty"`
	Title            string `json:"title"`
	Color            string `json:"color"`
	PublishedMobile  bool   `json:"published_mobile"`
	PublishedDesktop bool   `json:"published_desktop"`
	StartDate        string `json:"start_date"`
	EndDate          string `json:"end_date"`
}
