package transactionsparepartentities

type PurchaseRequestReferenceType struct {
	ReferencesTypeId   int    `gorm:"column:reference_type_id;size:50;not null;primaryKey;" json:"reference_type_id"`
	ReferencesTypeCode string `gorm:"column:reference_type_code;size:50" json:"reference_type_code"`
	ReferencesTypeName string `gorm:"column:reference_type_name;size:50;" json:"reference_type_name"`
}

func (*PurchaseRequestReferenceType) TableName() string {
	return "mtr_reference_type_purchase_request"
}
