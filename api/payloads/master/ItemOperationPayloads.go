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
	ItemId          int `json:"item_id"`
	ItemOperationId int `json:"item_operation_id"`
	LineTypeId      int `json:"line_type_id"`
	OperationId     int `json:"operation_id"`
	PackageId       int `json:"package_id"`
}

type MappingItemOperationResponse struct {
	ItemOperationId int    `json:"item_operation_id"`
	LineTypeId      int    `json:"line_type_id"`
	ItemId          int    `json:"item_id"`
	ItemCode        string `json:"item_code"`
	ItemName        string `json:"item_name"`
	OperationId     int    `json:"operation_id"`
	OperationCode   string `json:"operation_code"`
	OperationName   string `json:"operation_name"`
	PackageId       int    `json:"package_id"`
	PackageCode     string `json:"package_code"`
	PackageName     string `json:"package_name"`
}
