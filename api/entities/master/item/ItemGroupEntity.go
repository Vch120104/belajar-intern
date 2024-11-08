package masteritementities

var CreateItemGroupTable = "mtr_item_group"

type ItemGroup struct {
	IsActive        bool   `gorm:"column:is_active;default:true;not null" json:"is_active"`
	ItemGroupId     int    `gorm:"column:item_group_id;size:30;not null;primaryKey" json:"item_group_id"`
	ItemGroupCode   string `gorm:"column:item_group_code;not null;type:varchar(10);" json:"item_group_code"`
	ItemGroupName   string `gorm:"column:item_group_name;not null;type:varchar(100);" json:"item_group_name"`
	IsItemSparepart bool   `gorm:"column:is_item_sparepart;default:false; null" json:"is_item_sparepart"`
}

func (*ItemGroup) TableName() string {
	return CreateItemGroupTable
}
