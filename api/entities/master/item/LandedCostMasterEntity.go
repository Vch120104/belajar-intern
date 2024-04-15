package masteritementities

var CreateLandedCostTable = "mtr_landed_cost"

type LandedCost struct {
	IsActive               bool    `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	CompanyId              int     `gorm:"column:company_id;not null"        json:"company_id"`
	SupplierId             int     `gorm:"column:supplier_id;not null"        json:"supplier_id"`
	ShippingMethodId       int     `gorm:"column:shipping_method_id;not null"        json:"shipping_method_id"`
	LandedCostTypeId       int     `gorm:"column:landed_cost_type_id;not null"        json:"landed_cost_type_id"`
	LandedCostId           int     `gorm:"column:landed_cost_id;not null;primary_key"        json:"landed_cost_id"`
	LandedCostDescription  string  `gorm:"column:landed_cost_description;size:10;not null"        json:"landed_cost_description"`
	LandedCostMasterFactor float64 `gorm:"column:landed_cost_master_factor;not null"        json:"landed_cost_master_faction"`
}

func (*LandedCost) TableName() string {
	return CreateLandedCostTable
}
