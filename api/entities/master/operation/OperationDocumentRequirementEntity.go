package masteroperationentities

const TableNameOperationDocumentRequirement = "mtr_operation_document_requirement"

type OperationDocumentRequirement struct {
	IsActive                                bool   `gorm:"column:is_active;not null;default:true" json:"is_active"`
	OperationDocumentRequirementId          int    `gorm:"column:operation_document_requirement_id;size:30;not null;primaryKey" json:"operation_document_requirement_id"`
	OperationModelMappingId                 int    `gorm:"column:operation_model_mapping_id;size:30;not null;" json:"operation_model_mapping_id"`
	Line                                    int    `gorm:"column:line;size:30;null" json:"line"`
	OperationDocumentRequirementDescription string `gorm:"column:operation_document_requirement_description;null" json:"operation_document_requirement_description"`
}

func (*OperationDocumentRequirement) TableName() string {
	return TableNameOperationDocumentRequirement
}
