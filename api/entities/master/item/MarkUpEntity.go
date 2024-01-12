package masteritementities

var CreateMarkupMasterTable = "mtr_mark_up_master"

type MarkupMaster struct {
	IsActive                bool   `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	MarkupMasterId          int    `gorm:"column:markup_master_id;not null;primaryKey"        json:"markup_master_id"`
	MarkupMasterCode        string `gorm:"column:markup_master_code;unique;size:10;not null"        json:"markup_master_code"`
	MarkupMasterDescription string `gorm:"column:markup_master_description;size:20;not null"        json:"markup_master_description"`
}

func (*MarkupMaster) TableName() string {
	return CreateMarkupMasterTable
}
