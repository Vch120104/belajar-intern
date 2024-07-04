package masterpayloads

type ForecastMasterResponse struct {
	IsActive                   bool    `json:"is_active" parent_entity:"mtr_forecast_master"`
	ForecastMasterId           int     `json:"forecast_master_id" parent_entity:"mtr_forecast_master"`
	SupplierId                 int     `json:"supplier_id" parent_entity:"mtr_forecast_master"`
	MovingCodeId               int     `json:"moving_code_id" parent_entity:"mtr_forecast_master"`
	OrderTypeId                int     `json:"order_type_id" parent_entity:"mtr_forecast_master"`
	ForecastMasterLeadTime     float64 `json:"forecast_master_lead_time" parent_entity:"mtr_forecast_master"`
	ForecastMasterSafetyFactor float64 `json:"forecast_master_safety_factor" parent_entity:"mtr_forecast_master"`
	ForecastMasterOrderCycle   float64 `json:"forecast_master_order_cycle" parent_entity:"mtr_forecast_master"`
}

type ForecastMasterListResponse struct {
	IsActive                   bool    `json:"is_active" parent_entity:"mtr_forecast_master"`
	ForecastMasterId           int     `json:"forecast_master_id" parent_entity:"mtr_forecast_master" main_table:"mtr_forecast_master"`
	MovingCodeId               int     `json:"moving_code_id" parent_entity:"mtr_moving_code" references:"mtr_moving_code"`
	MovingCodeDescription      string  `json:"moving_code_description" parent_entity:"mtr_moving_code"`
	SupplierId                 int     `json:"supplier_id" parent_entity:"mtr_forecast_master"`
	OrderTypeId                int     `json:"order_type_id" parent_entity:"mtr_forecast_master"`
	ForecastMasterLeadTime     string  `json:"forecast_master_lead_time" parent_entity:"mtr_forecast_master"`
	ForecastMasterSafetyFactor string  `json:"forecast_master_safety_factor" parent_entity:"mtr_forecast_master"`
	ForecastMasterOrderCycle   float64 `json:"forecast_master_order_cycle" parent_entity:"mtr_forecast_master"`
}

type ForecastMasterResponseUpdate struct {
	IsActive                   bool    `json:"is_active" parent_entity:"mtr_forecast_master"`
	SupplierId                 int     `json:"supplier_id" parent_entity:"mtr_forecast_master"`
	MovingCodeId               int     `json:"moving_code_id" parent_entity:"mtr_forecast_master"`
	OrderTypeId                int     `json:"order_type_id" parent_entity:"mtr_forecast_master"`
	ForecastMasterLeadTime     float64 `json:"forecast_master_lead_time" parent_entity:"mtr_forecast_master"`
	ForecastMasterSafetyFactor float64 `json:"forecast_master_safety_factor" parent_entity:"mtr_forecast_master"`
	ForecastMasterOrderCycle   float64 `json:"forecast_master_order_cycle" parent_entity:"mtr_forecast_master"`
}

type OrderTypeResponse struct {
	OrderTypeId   int    `json:"order_type_id"`
	OrderTypeName string `json:"order_type_name"`
}

type SupplierResponse struct {
	SupplierId   int    `json:"supplier_id"`
	SupplierName string `json:"supplier_name"`
}
