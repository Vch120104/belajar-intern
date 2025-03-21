package masterentities

var CreatePackageMasterDetailTable = "mtr_package_master_detail"

type PackageMasterDetail struct {
	IsActive             bool                  `gorm:"column:is_active;not null;default:true" json:"is_active"`
	PackageDetailId      int                   `gorm:"column:package_detail_id;size:30;not null;primaryKey" json:"package_detail_id"`
	PackageId            int                   `gorm:"column:package_id;unique;size:30;not null" json:"package_id"`
	LineTypeId           int                   `gorm:"column:line_type_id;size:30;not null" json:"line_type_id"`
	ItemOperationId      int                   `gorm:"column:item_operation_id;size:30;not null" json:"item_operation_id"`
	FrtQuantity          float64               `gorm:"column:frt_quantity;not null" json:"frt_quantity"`
	TransactionTypeId    int                   `gorm:"column:transaction_type_id;size:30;not null" json:"transaction_type_id"`
	JobTypeId            int                   `gorm:"column:job_type_id;size:30;not null" json:"job_type_id"`
	MappingItemOperation *MappingItemOperation `gorm:"foreignKey:ItemOperationId;references:ItemOperationId"`
	Package              *PackageMaster        `gorm:"foreignKey:PackageId;references:PackageId"`
}

func (*PackageMasterDetail) TableName() string {
	return CreatePackageMasterDetailTable
}
