package masteroperationpayloads

type OperationCodeResponse struct {
	IsActive                bool   `json:"is_active"`
	OperationId             int32  `json:"operation_id"`
	OperationCode           string `json:"operation_code"`
	OperationName           string `json:"operation_name"`
	OperationUsingIncentive bool   `json:"operation_using_incentive"`
	OperationUsingActual    bool   `json:"operation_using_actual"`
}

type OperationCodeGetAll struct {
	OperationCode string `json:"operation_code"`
	OperationName string `json:"operation_name"`
	IsActive      bool   `json:"is_active"`
}

type OperationCodeSave struct {
	OperationCode           string `json:"operation_code"`
	OperationName           string `json:"operation_name"`
	IsActive                bool   `json:"is_active"`
	OperationUsingIncentive bool   `json:"operation_using_incentive"`
	OperationUsingActual    bool   `json:"operation_using_actual"`
}

type OperationCodeUpdate struct{
	OperationName           string `json:"operation_name"`
	OperationUsingIncentive bool   `json:"operation_using_incentive"`
	OperationUsingActual    bool   `json:"operation_using_actual"`
}
