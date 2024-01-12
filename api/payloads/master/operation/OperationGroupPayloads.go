package masteroperationpayloads

type OperationGroupResponse struct {
	IsActive                  bool   `json:"is_active"`
	OperationGroupId          int    `json:"operation_group_id"`
	OperationGroupCode        string `json:"operation_group_code"`
	OperationGroupDescription string `json:"operation_group_description"`
}

type ChangeStatusOperationGroupRequest struct {
	IsActive bool `json:"is_active"`
}
