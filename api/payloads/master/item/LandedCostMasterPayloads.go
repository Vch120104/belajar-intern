package masteritempayloads

type LandedCostMasterPayloads struct {
	IsActive               bool    `json:"is_active" parent_entity:"mtr_landed_cost"`
	LandedCostId           int     `json:"landed_cost_id" parent_entity:"mtr_landed_cost"`
	CompanyId              int     `json:"company_id" parent_entity:"mtr_landed_cost"`
	SupplierId             int     `json:"supplier_id" parent_entity:"mtr_landed_cost"`
	ShippingMethodId       int     `json:"shipping_method_id" parent_id:"mtr_landed_cost"`
	LandedCostTypeId       int     `json:"landed_cost_type_id" parent_entity:"mtr_landed_cost"`
	LandedCostDescription  string  `json:"landed_cost_description" parent_entity:"mtr_landed_cost"`
	LandedCostMasterFactor float64 `json:"landed_cost_master_factor" parent_entity:"mtr_landed_cost"`
}

type ShippingMethodResponse struct {
	ShippingMethodId   int    `json:"shipping_method_id"`
	ShippingMethodCode string `json:"shipping_method_code"`
}

type LandedCostTypeResponse struct {
	LandedCostTypeId          int    `json:"landed_cost_type_id"`
	LandedCostTypeCode        string `json:"landed_cost_type_code"`
	LandedCostTypeDescription string `json:"landed_cost_type_description"`
}
