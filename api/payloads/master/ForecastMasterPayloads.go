package masterpayloads

type ForecastMasterResponse struct {
	IsActive                   bool    `json:"is_active"`
	ForecastMasterId           int     `json:"forecast_master_id"`
	SupplierId                 int     `json:"supplier_id"`
	MovingCodeId               int     `json:"moving_code_id"`
	OrderTypeId                int     `json:"order_type_id"`
	ForecastMasterReadTime     float64 `json:"forecast_master_read_time"`
	ForecastMasterSafetyFactor float64 `json:"forecast_master_safety_factor"`
	ForecastMasterOrderCycle   float64 `json:"forecast_master_order_cycle"`
}
