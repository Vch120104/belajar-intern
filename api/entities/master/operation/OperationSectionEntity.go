package masteroperationentities

const TableNameOperationSection = "mtr_operation_section"

type OperationSection struct {
	IsActive                    bool   `gorm:"column:is_active;not null;default:true" json:"is_active"`
	OperationSectionId          int    `gorm:"column:operation_section_id;size:30;not null;primaryKey" json:"operation_section_id"`
	OperationSectionCode        string `gorm:"column:operation_section_code;uniqueIndex:idx_code_id;not null;type:char(3)" json:"operation_section_code"`
	OperationGroupId            int    `gorm:"column:operation_group_id;size:30;uniqueIndex:idx_code_id;not null;" json:"operation_group_id"`
	OperationGroup              *OperationGroup
	OperationSectionDescription string           `gorm:"column:operation_section_description;not null;size:50" json:"operation_section_description"`
	OperationEntries            OperationEntries `gorm:"foreignKey:OperationSectionId;references:OperationSectionId"`
	OperationKey                OperationKey     `gorm:"foreignKey:OperationSectionId;references:OperationSectionId"`
}

func (*OperationSection) TableName() string {
	return TableNameOperationSection
}
