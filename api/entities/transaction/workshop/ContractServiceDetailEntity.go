package transactionworkshopentities

const TableNameContractServiceDetail = "trx_contract_service_detail"

type ContractServiceDetail struct {
	ContractServicePackageDetailSystemNumber int     `gorm:"column:contract_service_package_detail_system_number;primary_key;size:30" json:"contract_service_package_detail_system_number"`
	ContractServiceSystemNumber              int     `gorm:"column:contract_service_system_number;size:30;not null" json:"contract_service"`
	ContractServiceLine                      string  `gorm:"column:contract_service_line;not null" json:"contract_service_line"`
	LineTypeId                               int     `gorm:"column:line_type_id;size:30;not null" json:"line_type_id"`
	ItemOperationId                          int     `gorm:"column:item_operation_id;size:30;not null" json:"item_operation_id"`
	Description                              string  `gorm:"column:description;not null" json:"description"`
	FrtQuantity                              float64 `gorm:"column:frt_quantity;not null" json:"frt_quantity"`
	ItemPrice                                float64 `gorm:"column:item_price;not null" json:"item_price"`
	ItemDiscountPercent                      float64 `gorm:"column:item_discount_percent;not null" json:"item_discount_percent"`
	ItemDiscountAmount                       float64 `gorm:"column:item_discount_amount;not null" json:"item_discount_amount"`
	PackageId                                int     `gorm:"column:package_id;size:30;not null" json:"package_id"`
	TotalUseFrtQuantity                      float64 `gorm:"column:total_use_frt_quantity;not null" json:"total_use_frt_quantity"`
}

func (*ContractServiceDetail) TableName() string {
	return TableNameContractServiceDetail
}
