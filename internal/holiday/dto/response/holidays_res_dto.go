package response

type GetHolidays4WebResDTO struct {
	HolidaysOrg []HolidaysByOrgRes `json:"holiday_organizations"`
}

type HolidaysByOrgRes struct {
	OrganizationName string          `json:"organization_name"`
	Holidays         []HolidayResDTO `json:"holidays"`
}
