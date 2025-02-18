package payloads

type UserDetail struct {
	UserId    int32  `json:"user_id"`
	Username  string `json:"username"`
	Authorize string `json:"authorized"`
	CompanyId string `json:"company_id"`
	Role      uint16 `json:"role"`
	IpAddress string `json:"ip_address"`
	Client    string `json:"client"`
}
