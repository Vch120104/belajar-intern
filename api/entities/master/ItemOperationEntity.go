package masterentities

const ItemOperationTableName = "mtr_item_operation"

type ItemOperation struct {
	ItemOperationId         int                   `gorm:"column:item_operation_id;not null;primaryKey;size:30" json:"item_operation_id"`
	ItemId                  int                   `gorm:"column:item_id;not null;size:30;uniqueIndex:item_cycle;default:0" json:"item_id"`
	OperationModelMappingId int                   `gorm:"column:operation_model_mapping_id;not null;uniqueIndex:item_cycle;default:0;size:30" json:"operation_model_mapping_id"`
	LineTypeId              int                   `gorm:"column:line_type_id;not null;size:30" json:"line_type_id"`
	PackageMasterDetail     *PackageMasterDetail  `gorm:"foreignKey:ItemOperationId;references:ItemOperationId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CampaignMasterDetail    *CampaignMasterDetail `gorm:"foreignKey:ItemOperationId;references:ItemOperationId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (*ItemOperation) TableName() string {
	return ItemOperationTableName
}
