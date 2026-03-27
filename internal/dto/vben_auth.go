package dto

type VbenLoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type VbenUserInfo struct {
	UserID   string   `json:"userId"`
	Username string   `json:"username"`
	RealName string   `json:"realName"`
	Roles    []string `json:"roles"`
}

type VbenPermissionCode = string
