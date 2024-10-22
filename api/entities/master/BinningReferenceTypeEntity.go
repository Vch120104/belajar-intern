package masterentities

// comGenVariable(VARIABLE like `BIN_REF_TYPE%`)
var BinningReferenceType = "mtr_binning_reference_type"

type BinningReferenceTypeMaster struct {
	BinningReferenceTypeId          int    `gorm:"column:binning_reference_type_id;;primaryKey"        json:"binning_reference_type_id"`
	BinningReferenceTypeCode        string `gorm:"column:binning_reference_type_code;"        json:"binning_reference_type_code"`
	BinningReferenceTypeDescription string `gorm:"column:binning_reference_type_description"        json:"binning_reference_type_description"`
}

func (*BinningReferenceTypeMaster) TableName() string {
	return BinningReferenceType
}
