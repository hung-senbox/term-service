package response

type GetTerms4WebResDTO struct {
	TermsOrg []TermsByOrgRes `json:"terms_org"`
}

type TermsByOrgRes struct {
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
