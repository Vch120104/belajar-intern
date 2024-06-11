package masteritementities

var CreateLandedCostTable = "mtr_landed_cost"

type LandedCost struct {
	IsActive         bool    `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	CompanyId        int     `gorm:"column:company_id;size:30;not null"        json:"company_id"`
	SupplierId       int     `gorm:"column:supplier_id;size:30;not null"        json:"supplier_id"`
	ShippingMethodId int     `gorm:"column:shipping_method_id;size:30;not null"        json:"shipping_method_id"`
	LandedCostTypeId int     `gorm:"column:landed_cost_type_id;size:30;not null"        json:"landed_cost_type_id"`
	LandedCostId     int     `gorm:"column:landed_cost_id;not null;size:30;primary_key"        json:"landed_cost_id"`
	LandedCostfactor float64 `gorm:"column:landed_cost_factor;not null" json:"landed_cost_factor"`
}

func (*LandedCost) TableName() string {
	return CreateLandedCostTable
}
