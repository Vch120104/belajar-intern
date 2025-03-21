package transactionworkshoppayloads

type ContractServiceDetailRequest struct {
	ContractServicePackageDetailSystemNumber int     `json:"contract_service_package_detail_system_number" parent_entity:"trx_contract_service_detail" main_table:"trx_contract_service_detail"`
	ContractServiceSystemNumber              int     `json:"contract_service_system_number" main_table:"trx_contract_service_detail"`
	ContractServiceLine                      string  `json:"contract_service_line" main_table:"trx_contract_service_detail"`
	LineTypeId                               int     `json:"line_type_id" main_table:"trx_contract_service_detail"`
	ItemOperationId                          int     `json:"item_operation_id" main_table:"trx_contract_service_detail"`
	Description                              string  `json:"description" main_table:"trx_contract_service_detail"`
	FrtQuantity                              float64 `json:"frt_quantity" main_table:"trx_contract_service_detail"`
	ItemPrice                                float64 `json:"item_price" main_table:"trx_contract_service_detail"`
	ItemDiscountPercent                      float64 `json:"item_discount_percent" main_table:"trx_contract_service_detail"`
	ItemDiscountAmount                       float64 `json:"item_discount_amount" main_table:"trx_contract_service_detail"`
	PackageId                                int     `json:"package_id" main_table:"trx_contract_service_detail"`
	TotalUseFrtQuantity                      float64 `json:"total_use_frt_quantity" main_table:"trx_contract_service_detail"`
}

type ContractServiceDetailPayloads struct {
	ContractServicePackageDetailSystemNumber int     `json:"contract_service_package_detail_system_number"`
	ContractServiceSystemNumber              int     `json:"contract_service_system_number"`
	ContractServiceLine                      string  `json:"contract_service_line"`
	LineTypeId                               int     `json:"line_type_id"`
	ItemOperationId                          int     `json:"item_operation_id"`
	Description                              string  `json:"description"`
	FrtQuantity                              float64 `json:"frt_quantity"`
	ItemPrice                                float64 `json:"item_price"`
	ItemDiscountPercent                      float64 `json:"item_discount_percent"`
	ItemDiscountAmount                       float64 `json:"item_discount_amount"`
	PackageId                                int     `json:"package_id"`
	TotalUseFrtQuantity                      float64 `json:"total_use_frt_quantity"`
}

type Operation struct {
	OperationCode string `json:"operation_code"`
	OperationName string `json:"operation_name"`
}

type ContractServiceIdResponse struct {
	ContractServicePackageDetailSystemNumber int     `json:"contract_service_package_detail_system_number"`
	ContractServiceSystemNumber              int     `json:"contract_service_system_number"`
	ContractServiceLine                      string  `json:"contract_service_line"`
	LineTypeId                               int     `json:"line_type_id"`
	ItemOperationId                          int     `json:"item_operation_id"`
	Description                              string  `json:"description"`
	FrtQuantity                              float64 `json:"frt_quantity"`
	ItemPrice                                float64 `json:"item_price"`
	ItemDiscountPercent                      float64 `json:"item_discount_percent"`
	ItemDiscountAmount                       float64 `json:"item_discount_amount"`
	PackageId                                int     `json:"package_id"`
	TotalUseFrtQuantity                      float64 `json:"total_use_frt_quantity"`
}

type LineTypeResponse struct {
	LineTypeId   int    `json:"line_type_id"`
	LineTypeCode string `json:"line_type_code"`
	LineTypeName string `json:"line_type_name"`
}
