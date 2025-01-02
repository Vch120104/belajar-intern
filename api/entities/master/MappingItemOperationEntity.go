package masterentities

import (
	masteritementities "after-sales/api/entities/master/item"
	masteroperationentities "after-sales/api/entities/master/operation"
)

const MappingItemOperationTableName = "mtr_mapping_item_operation"

type MappingItemOperation struct {
	ItemOperationId int                                    `gorm:"column:item_operation_id;size:30;not null;primaryKey" json:"item_operation_id"`
	LineTypeId      int                                    `gorm:"column:line_type_id;size:30;not null;" json:"line_type_id"`
	ItemId          int                                    `gorm:"column:item_id;size:30;" json:"item_id"`
	OperationId     int                                    `gorm:"column:operation_id;size:30;"  json:"operation_id"`
	PackageId       int                                    `gorm:"column:package_id;size:30;" json:"package_id"`
	Item            *masteritementities.Item               `gorm:"foreignKey:ItemId;references:ItemId"`
	Operation       *masteroperationentities.OperationCode `gorm:"foreignKey:OperationId;references:OperationId"`
	Package         *PackageMaster                         `gorm:"foreignKey:PackageId;references:PackageId"`
}

func (*MappingItemOperation) TableName() string {
	return MappingItemOperationTableName
}
