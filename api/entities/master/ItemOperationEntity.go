package masterentities

const ItemOperationTableName = "mtr_item_operation"

type ItemOperation struct {
	ItemOperationId             int `gorm:"column:item_operation_id;not null;primaryKey" json:"item_operation_id"`
	ItemOperationModelMappingId int `gorm:"column:item_operation_model_mapping_id;size:30;not null;uniqueIndex:item_cycle" json:"item_operation_model_mapping_id"`
	LineTypeId                  int `gorm:"column:line_type_id;not null;uniqueIndex:item_cycle" json:"line_type_id"`
}

func (*ItemOperation) TableName() string {
	return ItemOperationTableName
}
