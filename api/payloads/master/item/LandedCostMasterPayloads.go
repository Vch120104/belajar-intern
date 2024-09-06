package masteritempayloads

type LandedCostMasterPayloads struct {
	IsActive         bool    `json:"is_active"`
	LandedCostId     int     `json:"landed_cost_id"`
	CompanyId        int     `json:"company_id"`
	SupplierId       int     `json:"supplier_id"`
	ShippingMethodId int     `json:"shipping_method_id"`
	LandedCostTypeId int     `json:"landed_cost_type_id"`
	LandedCostFactor float64 `json:"landed_cost_factor"`
}

type ShippingMethodResponse struct {
	ShippingMethodId          int    `json:"shipping_method_id"`
	ShippingMethodCode        string `json:"shipping_method_code"`
	ShippingMethodDescription string `json:"shipping_method_description"`
}

type LandedCostTypeResponse struct {
	LandedCostTypeId          int    `json:"landed_cost_type_id"`
	LandedCostTypeCode        string `json:"landed_cost_type_code"`
	LandedCostTypeDescription string `json:"landed_cost_type_name"`
}

type LandedCostMasterRequest struct {
	IsActive         bool    `json:"is_acctive"`
	CompanyId        int     `json:"company_id"`
	SupplierId       int     `json:"supplier_id"`
	ShippingMethodId int     `json:"shipping_method_id"`
	LandedCostTypeId int     `json:"landed_cost_type_id"`
	LandedCostId     int     `json:"landed_cost_id"`
	LandedCostFactor float64 `json:"landed_cost_factor"`
}

type LandedCostMasterUpdateRequest struct {
	LandedCostfactor float64 `json:"landed_cost_factor"`
}
