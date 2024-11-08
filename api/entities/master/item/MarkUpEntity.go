package masteritementities

var CreateMarkupMasterTable = "mtr_markup_master"

type MarkupMaster struct {
	IsActive          bool   `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	MarkupMasterId    int    `gorm:"column:markup_master_id;size:30;not null;primaryKey"        json:"markup_master_id"`
	MarkupCode        string `gorm:"column:markup_code;unique;size:10;not null"        json:"markup_master_code"`
	MarkupDescription string `gorm:"column:markup_description;size:20;not null"        json:"markup_description"`
}

func (*MarkupMaster) TableName() string {
	return CreateMarkupMasterTable
}
