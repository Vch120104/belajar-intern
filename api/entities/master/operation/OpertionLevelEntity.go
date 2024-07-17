package masteroperationentities

const TableNameOperationLevel = "mtr_operation_level"

type OperationLevel struct {
	OperationDocumentRequirementId          int    `gorm:"column:operation_document_requirement_id;not null;primaryKey;size:30" json:"operation_document_requirement_id"`
	OperationModelMappingId                 int    `gorm:"column:operation_model_mapping_id;not null;size:30" json:"operation_model_mapping_id"`
	Line                                    int    `gorm:"column:line;not null;size:30;unique" json:"line"`
	OperationDocumentRequirementDescription string `gorm:"column:operation_document_requirement_document;not null;size:50" json:"operation_document_requirement_document"`
}

func (*OperationLevel) TableName() string {
	return TableNameOperationLevel
}
