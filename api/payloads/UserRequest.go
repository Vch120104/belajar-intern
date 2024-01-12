package payloads

type UserDetail struct {
	UserId      int32  `json:"user_id"`
	Username    string `json:"username"`
	Authorize   string `json:"authorized"`
	CompanyCode string `json:"company_code"`
	Role        uint16 `json:"role"`
	IpAddress   string `json:"ip_address"`
	Client      string `json:"client"`
}
