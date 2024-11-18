package masterpayloads

type GetOrderTypeResponse struct {
	IsActive      bool   `json:"is_active"`
	OrderTypeId   int    `json:"order_type_id"`
	OrderTypeCode string `json:"order_type_code"`
	OrderTypeName string `json:"order_type_name"`
}

type OrderTypeSaveRequest struct {
	OrderTypeCode string `json:"order_type_code" validate:"required"`
	OrderTypeName string `json:"order_type_name" validate:"required"`
}

type OrderTypeUpdateRequest struct {
	OrderTypeCode string `json:"order_type_code" validate:"required"`
	OrderTypeName string `json:"order_type_name" validate:"required"`
}
