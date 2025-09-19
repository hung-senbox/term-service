package response

type GetHolidays4WebResDTO struct {
	HolidaysOrg []HolidaysByOrgRes `json:"holidays_org"`
}

type HolidaysByOrgRes struct {
	OrganizationName string          `json:"organization_name"`
	Holidays         []HolidayResDTO `json:"holidays"`
}
