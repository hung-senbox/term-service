package response

type ListTermsOrgResDTO struct {
	TermsOrg []TemsByOrgRes `json:"terms_org"`
}

type TemsByOrgRes struct {
	OrganizationName string       `json:"organization_name"`
	Terms            []TermResDTO `json:"terms"`
}

/////// LIST TERMS

type ListTermsResDTO struct {
	Terms []TermResDTO `json:"terms"`
}

type GetTerms4AppResDTO struct {
	Terms []CurrentTermResDTO `json:"terms"`
}
