package masterentities

var CreateMovingCodeTable = "mtr_moving_code"

type MovingCode struct {
	IsActive              bool    `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	MovingCodeId          int     `gorm:"column:moving_code_id;size:30;not null;primaryKey"        json:"moving_code_id"`
	CompanyId             int     `gorm:"column:company_id;size:30;not null"        json:"company_id"`
	MovingCodeDescription string  `gorm:"column:moving_code_description;size:40;not null" json:"moving_code_description"`
	MinimumQuantityDemand float64 `gorm:"column:minimum_quantity_demand;null" json:"minimum_quantity_demand"`
	Priority              float64 `gorm:"column:priority;null" json:"priority"`
	AgingMonthFrom        float64 `gorm:"column:aging_month_from;null" json:"aging_month_from"`
	AgingMonthTo          float64 `gorm:"column:aging_month_to;null" json:"aging_month_to"`
	DemandExistMonthFrom  float64 `gorm:"column:demand_exist_month_from;null" json:"demand_exist_month_from"`
	DemandExistMonthTo    float64 `gorm:"column:demand_exist_month_to;null" json:"demand_exist_month_to"`
	LastMovingMonthFrom   float64 `gorm:"column:last_moving_month_from;null" json:"last_moving_month_from"`
	LastMovingMonthTo     float64 `gorm:"column:last_moving_month_to;null" json:"last_moving_month_to"`
	Remark                string  `gorm:"column:remark;null" json:"remark"`
}

func (*MovingCode) TableName() string {
	return CreateMovingCodeTable
}
