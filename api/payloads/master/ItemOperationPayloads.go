package masterpayloads

type ItemOperationGet struct {
	ItemOperationId         int    `json:"item_operation_id"`
	ItemId                  int    `json:"item_id"`
	ItemName                string `json:"item_name"`
	OperationModelMappingId int    `json:"operation_model_mapping_id"`
	OperationName           string `json:"operation_name"`
	LineTypeId              int    `json:"line_type_id"`
}

type ItemOperationPost struct {
	ItemId                  int `json:"item_id"`
	OperationModelMappingId int `json:"operation_model_mapping_id"`
	LineTypeId              int `json:"line_type_id"`
}
