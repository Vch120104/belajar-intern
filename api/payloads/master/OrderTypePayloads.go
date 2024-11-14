package masterpayloads

type GetOrderTypeResponse struct {
	IsActive      bool   `json:"is_active"`
	OrderTypeId   int    `json:"order_type_id"`
	OrderTypeCode string `json:"order_type_code"`
	OrderTypename string `json:"order_type_name"`
}
