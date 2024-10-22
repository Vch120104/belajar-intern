package transactionworkshopentities

const TableNameContractServiceOperationDetail = "trx_contract_service_operation_detail"

type ContractServiceOperationDetail struct {
	ContractServicePackageDetailSystemNumber int     `gorm:"column:contract_service_package_detail_system_number;primary_key;size:30" json:"contract_service_package_detail_system_number"`
	ContractServiceSystemNumber              int     `gorm:"column:contract_service_system_number;size:30;not null" json:"contract_service"`
	ContractServiceLine                      string  `gorm:"column:contract_service_line;not null" json:"contract_service_line"`
	LineTypeId                               int     `gorm:"column:line_type_id;size:30;not null" json:"line_type_id"`
	OperationId                              int     `gorm:"column:operation_id;size:30;not null" json:"operation_id"`
	Description                              string  `gorm:"column:description;not null" json:"description"`
	FrtQuantity                              float64 `gorm:"column:frt_quantity;not null" json:"frt_quantity"`
	OperationPrice                           float64 `gorm:"column:operation_price;not null" json:"operation_price"`
	OperationDiscountPercent                 float64 `gorm:"column:operation_discount_percent;not null" json:"operation_discount_percent"`
	OperationDiscountAmount                  float64 `gorm:"column:operation_discount_amount;not null" json:"operation_discount_amount"`
	PackageId                                int     `gorm:"column:package_id;size:30;not null" json:"package_id"`
	TotalUseFrtQuantity                      float64 `gorm:"column:total_use_frt_quantity;not null" json:"total_use_frt_quantity"`
}

func (*ContractServiceOperationDetail) TableName() string {
	return TableNameContractServiceOperationDetail
}
