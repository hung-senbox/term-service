package request

type GetAvatarUrlRequest struct {
	OwnerID   string `json:"owner_id"`
	OwnerRole string `json:"owner_role"`
}
