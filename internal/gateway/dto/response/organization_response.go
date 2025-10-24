package response

type OrganizationInfo struct {
	ID               string                      `json:"id"`
	OrganizationName string                      `json:"organization_name"`
	Avatar           string                      `json:"avatar"`
	AvatarURL        string                      `json:"avatar_url"`
	Address          string                      `json:"address"`
	Description      string                      `json:"description"`
	Managers         []GetOrgManagerInfoResponse `json:"managers"`
}

type GetOrgManagerInfoResponse struct {
	UserID       string `json:"user_id"`
	UserNickName string `json:"user_nick_name"`
	IsManager    bool   `json:"is_manager"`
}
