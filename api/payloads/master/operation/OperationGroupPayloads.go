package masteroperationpayloads

type OperationGroupResponse struct {
	IsActive                  bool   `json:"is_active" `
	OperationGroupId          int    `json:"operation_group_id"`
	OperationGroupCode        string `json:"operation_group_code" validate:"required,min=1,max=2"`
	OperationGroupDescription string `json:"operation_group_description" validate:"required,max=50"`
}

type OperationGroupDropDownResponse struct {
	IsActive           bool   `json:"is_active" `
	OperationGroupId   int    `json:"operation_group_id"`
	OperationGroupCode string `json:"operation_group_code_description"`
}

type ChangeStatusOperationGroupRequest struct {
	IsActive bool `json:"is_active"`
}
