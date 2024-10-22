package masterentities

var BinningType = "mtr_binning_type"

type BinningTypeMaster struct {
	BinningTypeId          int    `gorm:"column:binning_type_id;not null;primaryKey" json:"binning_type_id"`
	BinningTypeCode        string `gorm:"column:binning_type_code;not null" json:"binning_type_code"`
	BinningTypeDescription string `gorm:"column:binning_type_description;not null" json:"binning_type_description"`
}

func (*BinningTypeMaster) TableName() string {
	return BinningType
}
