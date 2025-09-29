package request

type UploadTermItem struct {
	ID               string `json:"id,omitempty"`
	Title            string `json:"title" bninding:"required"`
	Color            string `json:"color" binding:"required"`
	PublishedMobile  bool   `json:"published_mobile"`
	PublishedDesktop bool   `json:"published_desktop"`
	PublishedTeacher bool   `json:"published_teacher"`
	PublishedParent  bool   `json:"published_parent"`
	StartDate        string `json:"start_date" bninding:"required"`
	EndDate          string `json:"end_date" binding:"required"`
}

type UploadTermRequest struct {
	LanguageID uint             `json:"language_id" binding:"required"`
	Word       string           `json:"word" binding:"required"`
	Terms      []UploadTermItem `json:"terms" binding:"required"`
}
