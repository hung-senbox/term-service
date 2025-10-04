package response

import "term-service/internal/gateway/dto"

type GetTerms4WebResDTO struct {
	TermsOrg []TermsByOrgRes `json:"term_organizations"`
}

type TermsByOrgRes struct {
	MessageLanguages []dto.MessageLanguageResponse `json:"message_languages"`
	OrganizationName string                        `json:"organization_name"`
	Terms            []TermResDTO                  `json:"terms"`
}

/////// LIST TERMS

type ListTermsResDTO struct {
	Terms []TermResDTO `json:"terms"`
}

type GetTerms4AppResDTO struct {
	Terms []CurrentTermResDTO `json:"terms"`
}

type Term4GwResponse struct {
	ID        string `json:"id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
