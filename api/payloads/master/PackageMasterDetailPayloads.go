
package masterpayloads

type PackageMasterDetailItem struct {
	IsActive                   bool    `json:"is_active"`
	PackageDetailItemId        int     `json:"package_detail_item_id"`
	PackageId                  int     `json:"package_id"`
	LineTypeId                 int     `json:"line_type_id"`
	ItemId                     int     `json:"item_id"`
	FrtQuantity                float64 `json:"frt_quantity"`
	WorkorderTransactionTypeId int     `json:"workorder_transaction_type_id"`
	JobTitleId                 int     `json:"job_title_id"`
}

type PackageMasterDetailOperation struct {
	IsActive                   bool    `json:"is_active"`
	PackageDetailItemId        int     `json:"package_detail_item_id"`
	PackageId                  int     `json:"package_id"`
	LineTypeId                 int     `json:"line_type_id"`
	OperationId                int     `json:"operation_id"`
	FrtQuantity                float64 `json:"frt_quantity"`
	WorkorderTransactionTypeId int     `json:"workorder_transaction_type_id"`
	JobTypeId                  int     `json:"job_type_id"`
}

type PackageMasterCombinedData struct {
	IsActive                   bool    `json:"is_active"`
	PackageDetailItemId        int     `json:"package_detail_item_id"`
	PackageId                  int     `json:"package_id"`
	LineTypeId                 int     `json:"line_type_id"`
	ItemOperationCode          string  `json:"item_operation_code"`
	ItemOperationName          string  `json:"item_operation_name"`
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
