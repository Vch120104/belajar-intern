package masterpackagemasterentity

import (
	masterentities "after-sales/api/entities/master"
	masteroperationentities "after-sales/api/entities/master/operation"
)

var CreatePackageMasterDetailOperationTable = "mtr_package_master_detail_operation"

type PackageMasterDetailOperation struct {
	IsActive                 bool `gorm:"collumn:is_active;not null; default:true" json:"is_actve"`
	PackageDetailOperationId int  `gorm:"column:package_detail_operation_id;size:30;not null;primary_key" json:"package_detail_operation_id"`
	PackageId                int  `gorm:"column:package_id;size:30;not null" json:"package_id"`
	Package                  *masterentities.PackageMaster
	LineTypeId               int `gorm:"column:line_type_id;size:30;not null" json:"line_type_id"`
	OperationId              int `gorm:"column:operation_id;size:30;not null" json:"operation_id"`
	Operation                *masteroperationentities.OperationModelMapping
	FrtQuantity              float64 `gorm:"column:frt_quantity;not null;" json:"frt_quantity"`
	Sequence                 int     `gorm:"column:sequesnce;size:30;not null;default:0" json:"sequence"`
	TransactionTypeId        int     `gorm:"column:workorder_transaction_type_id;size:30;not null;default:0" json:"workorder_transaction_type_id"`
	JobTypeId                int     `gorm:"column:job_type_id;size:30;not null;default:0" json:"job_type_id"`
}

func (*PackageMasterDetailOperation) TableName() string {
	return CreatePackageMasterDetailOperationTable
}