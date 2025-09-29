package request

type UploadHolidayItem struct {
	ID               string `json:"id,omitempty"`
	Title            string `json:"title" binding:"required"`
	Color            string `json:"color" binding:"required"`
	PublishedMobile  bool   `json:"published_mobile"`
	PublishedDesktop bool   `json:"published_desktop"`
	StartDate        string `json:"start_date" binding:"required"`
	EndDate          string `json:"end_date" binding:"required"`
}

type UploadHolidayRequest struct {
	LanguageID uint                `json:"language_id" binding:"required"`
	DeleteIds  []string            `json:"delete_ids"`
	Holidays   []UploadHolidayItem `json:"holidays"`
}
