package masteritementities

var CreateItemPackageDetailTable = "mtr_item_package_detail"

type ItemPackageDetail struct {
	ItemPackageDetailId int     `gorm:"column:item_package_detail_id;not null;primaryKey;size:30"        json:"item_package_detail_id"`
	IsActive            bool    `gorm:"column:is_active;not null"        json:"is_active"`
	ItemPackageId       int     `gorm:"column:item_package_id;not null;size:30"        json:"item_package_id"`
	ItemId              int     `gorm:"column:item_id;not null;size:30;unique"        json:"item_id"`
	Item                Item    `gorm:"foreignKey:ItemId;references:ItemId" json:"item"`
	Quantity            float64 `gorm:"column:quantity;not null"        json:"quantity"`
}

func (*ItemPackageDetail) TableName() string {
	return CreateItemPackageDetailTable
}
