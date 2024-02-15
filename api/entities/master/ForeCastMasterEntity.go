package masterentities

var CreateForecastMasterTable = "mtr_forecast_master"

type ForecastMaster struct {
	IsActive                   bool `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	ForecastMasterId           int  `gorm:"column:forecast_master_id;not null;primaryKey"        json:"forecast_master_id"`
	CompanyId                  int  `gorm:"column:company_id;not null"        json:"company_id"`
	SupplierId                 int  `gorm:"column:supplier_id;not null" json:"supplier_id"`
	MovingCodeId               int  `gorm:"column:moving_code_id;not null"      json:"moving_code_id"`
	MovingCode                 MovingCode
	OrderTypeId                int     `gorm:"column:order_type_id;not null" json:"order_type_id"`
	ForecastMasterLeadTime     float64 `gorm:"column:forecast_master_lead_time;not null" json:"forecast_master_lead_time"`
	ForecastMasterSafetyFactor float64 `gorm:"column:forecast_master_safety_factor;not null" json:"forecast_master_safety_factor"`
	ForecastMasterOrderCycle   float64 `gorm:"column:forecast_master_order_cycle;not null" json:"forecast_master_order_cycle"`
}

func (*ForecastMaster) TableName() string {
	return CreateForecastMasterTable
}
