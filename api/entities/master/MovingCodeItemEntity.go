package masterentities

import "time"

const MovingCodeItem = "mtr_moving_code_item"

type MovingItemCode struct {
	IsActive         bool       `gorm:"column:is_active;not null;default:true" json:"is_active"`
	MovingCodeItemId int        `gorm:"column:moving_code_item_id;size:30;not null;primaryKey" json:"moving_code_item_id"`
	CompanyId        int        `gorm:"column:company_id;size:30;not null;uniqueindex:item_cycle" json:"company_id"`
	ProcessDate      time.Time  `gorm:"column:process_date;not null;uniqueindex:item_cycle" json:"process_date"`
	ItemId           int        `gorm:"column:item_id;size:30;not null;uniqueindex:item_cycle" json:"item_id"`
	MovingCodeId     int        `gorm:"column:moving_code_id;size:30;null;" json:"moving_code_id"`
	MovingCode       MovingCode `gorm:"foreignKey:MovingCodeId;references:MovingCodeId" json:"moving_code"`
}

func (*MovingItemCode) TableName() string {
	return MovingCodeItem
}
