package masterpayloads

type AgreementResponse struct {
	IsActive              bool    `json:"is_active" parent_entity:"mtr_forecast_master"`
	AgreementId           int     `json:"forecast_master_id" parent_entity:"mtr_forecast_master"`
	CustomerId            int     `json:"customer_id" parent_entity:"mtr_forecast_master"`
	MovingCodeId          int     `json:"moving_code_id" parent_entity:"mtr_forecast_master"`
	OrderTypeId           int     `json:"order_type_id" parent_entity:"mtr_forecast_master"`
	AgreementLeadTime     float64 `json:"forecast_master_lead_time" parent_entity:"mtr_forecast_master"`
	AgreementSafetyFactor float64 `json:"forecast_master_safety_factor" parent_entity:"mtr_forecast_master"`
	AgreementOrderCycle   float64 `json:"forecast_master_order_cycle" parent_entity:"mtr_forecast_master"`
}

type AgreementListResponse struct {
	IsActive              bool    `json:"is_active" parent_entity:"mtr_forecast_master"`
	AgreementId           int     `json:"forecast_master_id" parent_entity:"mtr_forecast_master" main_table:"mtr_forecast_master"`
	MovingCodeId          int     `json:"moving_code_id" parent_entity:"mtr_moving_code" references:"mtr_moving_code"`
	MovingCodeDescription string  `json:"moving_code_description" parent_entity:"mtr_moving_code"`
	SupplierId            int     `json:"supplier_id" parent_entity:"mtr_forecast_master"`
	OrderTypeId           int     `json:"order_type_id" parent_entity:"mtr_forecast_master"`
	AgreementLeadTime     string  `json:"forecast_master_lead_time" parent_entity:"mtr_forecast_master"`
	AgreementSafetyFactor string  `json:"forecast_master_safety_factor" parent_entity:"mtr_forecast_master"`
	AgreementOrderCycle   float64 `json:"forecast_master_order_cycle" parent_entity:"mtr_forecast_master"`
}
