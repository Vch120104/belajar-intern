package masterpayloads

type ItemOperationGet struct {
	ItemOperationId int    `json:"item_operation_id"`
	ItemId          int    `json:"item_id"`
	ItemName        string `json:"item_name"`
	OperationId     int    `json:"operation_id"`
	OperationName   string `json:"operation_name"`
	LineTypeId      int    `json:"line_type_id"`
}

type ItemOperationPost struct{
	ItemId          int    `json:"item_id"`
	OperationId     int    `json:"operation_id"`
	LineTypeId      int    `json:"line_type_id"`
}
