package dto

type StudentResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	StudentName    string `json:"student_name"`
}
