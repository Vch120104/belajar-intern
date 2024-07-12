package masterpackagemasterentity

import (
	masterentities "after-sales/api/entities/master"
	masteritementities "after-sales/api/entities/master/item"
)

var CreatePackageMasterDetailItemTable = "mtr_package_master_detail_item"

type PackageMasterDetailItem struct {
	IsActive                   bool `gorm:"column:is_active;not null;default:true" json:"is_active"`
	PackageDetailItemId        int  `gorm:"column:package_detail_item_id;size:30;not null;primary_key" json:"package_detail_item_id"`
	PackageId                  int  `gorm:"column:package_id;size:30;not null" json:"package_id"`
	Package                    *masterentities.PackageMaster
	LineTypeId                 int `gorm:"column:line_type_id;size:30;not null" json:"line_type_id"`
	ItemId                     int `gorm:"column:item_id;size:30;not null" json:"item_id"`
	Item                       *masteritementities.Item
	FrtQuantity                float64 `gorm:"column:frt_quantity;not null" json:"frt_quantity"`
	WorkorderTransactionTypeId int     `gorm:"column:workorder_transaction_type_id;size:30;not null" json:"workorder_transaction_type_id"`
	JobTypeId                  int     `gorm:"column:job_type_id;size:30;not null" json:"job_type_id"`
}

func (*PackageMasterDetailItem) TableName() string {
	return CreatePackageMasterDetailItemTable
}