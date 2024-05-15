package masterpayloads

type PackageMasterDetailItem struct {
	IsActive                   bool    `json:"is_active" parent_entity:"mtr_package_master_detail_item"`
	PackageDetailItemId        int     `json:"package_detail_item_id" parent_entity:"mtr_package_master_detail_item"`
	PackageId                  int     `json:"package_id" parent_entity:"mtr_package_master_detail_item"`
	LineTypeId                 int     `json:"line_type_id" parent_entity:"mtr_package_master_detail_item"`
	ItemId                     int     `json:"item_id" parent_entity:"mtr_package_master_detail_item"`
	ItemCode                   string  `json:"item_code" parent_entity:"mtr_item"`
	ItemName                   string  `json:"item_name" parent_entity:"mtr_item"`
	FrtQuantity                float64 `json:"frt_quantity" parent_entity:"mtr_package_master_detail_item"`
	WorkorderTransactionTypeId int     `json:"workorder_transaction_type_id" parent_entity:"mtr_package_master_detail_item"`
	JobTypeId                  int     `json:"job_type_id" parent_entity:"mtr_package_master_detail_item"`
}

type PackageMasterDetailOperation struct {
	IsActive                   bool    `json:"is_active" parent_entity:"mtr_package_master_detail_operation"`
	PackageDetailOperationId   int     `json:"package_detail_operation_id" parent_entity:"mtr_package_master_detail_operation"`
	PackageId                  int     `json:"package_id" parent_entity:"mtr_package_master_detail_operation"`
	LineTypeId                 int     `json:"line_type_id" parent_entity:"mtr_package_master_detail_operation"`
	OperationId                int     `json:"operation_id" parent_entity:"mtr_package_master_detail_operation"`
	OperationName              string  `json:"operation_name" parent_entity:"mtr_operation_code"`
	OperationCode              string  `json:"operation_code" parent_entity:"mtr_operation_code"`
	FrtQuantity                float64 `json:"frt_quantity" parent_entity:"mtr_package_master_detail_operation"`
	WorkorderTransactionTypeId int     `json:"workorder_transaction_type_id" parent_entity:"mtr_package_master_detail_operation"`
	JobTypeId                  int     `json:"job_type_id" parent_entity:"mtr_package_master_detail_operation"`
}

type PackageMasterCombinedData struct {
	IsActive                   bool    `json:"is_active"`
	PackageDetailItemId        int     `json:"package_detail_item_id"`
	PackageId                  int     `json:"package_id"`
	LineTypeId                 int     `json:"line_type_id"`
	ItemId                     int     `json:"item_id"`
	OperationId                int     `json:"operation_id"`
	FrtQuantity                float64 `json:"frt_quantity"`
	WorkorderTransactionTypeId int     `json:"workorder_transaction_type_id"`
	JobTypeId                  int     `json:"job_type_id"`
}

type PackageMasterDetailOperationBodyshop struct {
	IsActive                 bool `json:"is_active"`
	PackageDetailOperationId int  `json:"package_detail_operation_id"`
	PackageId                int  `json:"package_id"`
	LineTypeId               int  `json:"line_type_id"`
	OperationId              int  `json:"operation_id"`
	Sequence                 int  `json:"sequence"`
}

type PackageMasterDetailWorkshop struct {
	IsActive                   bool    `json:"is_active"`
	PackageDetailItemId        int     `json:"package_detail_item_id"`
	PackageId                  int     `json:"package_id"`
	LineTypeId                 int     `json:"line_type_id"`
	ItemOperationId            int     `json:"item_operation_id"`
	FrtQuantity                float64 `json:"frt_quantity"`
	WorkorderTransactionTypeId int     `json:"workorder_transaction_type_id"`
	JobTypeId                  int     `json:"job_type_id"`
}

type LineTypeCode struct {
	LineTypeId   int    `json:"line_type_id"`
	LineTypeCode string `json:"line_type_code"`
}
