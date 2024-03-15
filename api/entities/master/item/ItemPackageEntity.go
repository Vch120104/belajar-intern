package masteritementities

var CreateItemPackageTable = "mtr_item_package"

type ItemPackage struct {
	ItemPackageId   int    `gorm:"column:item_package_id;not null;primaryKey;size:30"        json:"item_package_id"`
	IsActive        bool   `gorm:"column:is_active;not null"        json:"is_active"`
	ItemGroupId     int    `gorm:"column:item_group_id;not null;size:30"        json:"item_group_id"`
	ItemPackageCode string `gorm:"column:item_package_code;not null"        json:"item_package_code"`
	ItemPackageName string `gorm:"column:item_package_name;not null"        json:"item_package_name"`
	ItemPackageSet  bool   `gorm:"column:item_package_set;not null"        json:"item_package_set"`
	Description     string `gorm:"column:description;null"        json:"description"`
}

func (*ItemPackage) TableName() string {
	return CreateItemPackageTable
}
