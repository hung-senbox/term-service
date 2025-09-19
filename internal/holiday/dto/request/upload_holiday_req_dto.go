package request

type UploadHolidayItem struct {
	ID               string `json:"id,omitempty"`
	Title            string `json:"title"`
	Color            string `json:"color"`
	PublishedMobile  bool   `json:"published_mobile"`
	PublishedDesktop bool   `json:"published_desktop"`
	StartDate        string `json:"start_date"`
	EndDate          string `json:"end_date"`
}

type UploadHolidayRequest struct {
	DeleteIds []string            `json:"delete_ids"`
	Holidays  []UploadHolidayItem `json:"holidays"`
}
