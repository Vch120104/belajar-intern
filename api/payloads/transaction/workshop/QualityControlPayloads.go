package transactionworkshoppayloads

type QualityControlResponse struct {
	WorkOrderDocumentNumber string `json:"work_order_document_number"`
	WorkOrderDate           string `json:"work_order_date"`
	ModelId                 int    `json:"model_id"`
	ModelName               string `json:"model_name"`
	ModelCode               string `json:"model_code"`
	VariantId               int    `json:"variant_id"`
	VariantName             string `json:"variant_name"`
	VarianCode              string `json:"varian_code"`
	VehicleId               int    `json:"vehicle_id"`
	VehicleCode             string `json:"vehicle_chassis_number"`
	VehicleTnkb             string `json:"vehicle_registration_certificate_tnkb"`
	CustomerId              int    `json:"customer_id"`
	CustomerName            string `json:"customer_name"`
	WorkOrderSystemNumber   int    `json:"work_order_system_number"`
	BrandId                 int    `json:"brand_id"`
	BrandCode               string `json:"brand_code"`
}

type QualityControlRequest struct {
	WorkOrderSystemNumber int `json:"work_order_system_number" parent_entity:"trx_work_order" main_table:"trx_work_order"`
	BrandId               int `json:"brand_id" parent_entity:"trx_work_order"`
	ModelId               int `json:"model_id" parent_entity:"trx_work_order"`
	VariantId             int `json:"variant_id" parent_entity:"trx_work_order"`
	VehicleId             int `json:"vehicle_id" parent_entity:"trx_work_order"`
	CustomerId            int `json:"customer_id" parent_entity:"trx_work_order"`
}

type QualityControlRequestId struct {
	WorkOrderSystemNumber int `json:"work_order_system_number" parent_entity:"trx_work_order" main_table:"trx_work_order"`
	BrandId               int `json:"brand_id" parent_entity:"trx_work_order"`
	ModelId               int `json:"model_id" parent_entity:"trx_work_order"`
	VariantId             int `json:"variant_id" parent_entity:"trx_work_order"`
	VehicleId             int `json:"vehicle_id" parent_entity:"trx_work_order"`
	CustomerId            int `json:"customer_id" parent_entity:"trx_work_order"`
}

type QualityControlIdResponse struct {
	WorkOrderDocumentNumber string                        `json:"work_order_document_number"`
	WorkOrderDate           string                        `json:"work_order_date"`
	BrandName               string                        `json:"brand_name"`
	ModelName               string                        `json:"model_name"`
	VariantName             string                        `json:"variant_name"`
	ColourName              string                        `json:"colour_name"`
	VehicleCode             string                        `json:"vehicle_chassis_number"`
	EngineCode              string                        `json:"vehicle_engine_number"`
	LastMilage              int                           `json:"last_milage"`
	CurrentMilage           int                           `json:"current_milage"`
	VehicleTnkb             string                        `json:"vehicle_registration_certificate_tnkb"`
	CustomerName            string                        `json:"customer_name"`
	Address0                string                        `json:"address_0"`
	Address1                string                        `json:"address_1"`
	Address2                string                        `json:"address_2"`
	RTRW                    string                        `json:"rt_rw"`
	Phone                   string                        `json:"phone"`
	ForemanName             string                        `json:"foreman_name"`
	ServiceAdvisorName      string                        `json:"service_advisor_name"`
	OrderDateTime           string                        `json:"order_date_time"`
	EstimatedFinished       string                        `json:"estimated_finished"`
	QualityControlDetails   QualityControlDetailsResponse `json:"quality_control_details"`
}

type QualityControlDetailsResponse struct {
	Page       int                            `json:"page"`
	Limit      int                            `json:"limit"`
	TotalPages int                            `json:"total_pages"`
	TotalRows  int                            `json:"total_rows"`
	Data       []QualityControlDetailResponse `json:"data"`
}

type QualityControlDetailResponse struct {
	OperationItemId   int     `json:"operation_item_id"`
	OperationItemCode string  `json:"operation_item_code"`
	OperationItemName string  `json:"operation_item_name"`
	Frt               float64 `json:"frt"`
	ServiceStatusId   int     `json:"service_status_id"`
	ServiceStatusName string  `json:"service_status_name"`
	TechnicianId      int     `json:"technician_id"`
	TechnicianName    string  `json:"technician_name"`
	TechnicianCode    string  `json:"technician_code"`
}

type QualityControlUpdateResponse struct {
	WorkOrderSystemNumber int    `json:"work_order_system_number"`
	WorkOrderDetailId     int    `json:"work_order_detail_id"`
	WorkOrderStatusId     int    `json:"work_order_status_id"`
	WorkOrderStatusName   string `json:"work_order_status_name"`
}

type QualityControlReorder struct {
	ExtraTime float64 `json:"quality_control_extra_frt"`
	Reason    string  `json:"quality_control_extra_reason"`
}
