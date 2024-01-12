package masteroperationpayloads

type OperationCodeResponse struct {
	IsActive                bool   `json:"is_active"`
	OperationId             int32  `json:"operation_id"`
	OperationCode           string `json:"operation_code"`
	OperationName           string `json:"operation_name"`
	OperationGroupId        int32  `json:"operation_group_id"`
	OperationSectionId      int32  `json:"operation_section_id"`
	OperationKeyId          int32  `json:"operation_key_id"`
	OperationEntriesId      int32  `json:"operation_entries_id"`
	OperationUsingIncentive bool   `json:"operation_using_incentive"`
	OperationUsingActual    bool   `json:"operation_using_actual"`
}

type OperationCodeRequest struct {
	OperationCode           string `json:"operation_code"`
	OperationName           string `json:"operation_name"`
	OperationGroupId        int32  `json:"operation_group_id"`
	OperationSectionId      int32  `json:"operation_section_id"`
	OperationKeyId          int32  `json:"operation_key_id"`
	OperationEntriesId      int32  `json:"operation_entries_id"`
	OperationUsingIncentive bool   `json:"operation_using_incentive"`
	OperationUsingActual    bool   `json:"operation_using_actual"`
}
