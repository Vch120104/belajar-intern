package transactionworkshoppayloads

import "time"

type VehicleHistoryResponses struct {
	WorkOrderSystemNumber      int       `json:"work_order_system_number" parent_entity:"trx_work_order"`
	WorkOrderDocumentNumber    string    `json:"work_order_document_number" parent_entity:"trx_work_order" main_table:"trx_work_order"`
	WorkOrderDate              time.Time `json:"service_request_date" parent_entity:"trx_work_order"`
	WorkOrderStatusDescription string    `json:"work_order_status_desc" parent_entity:"mtr_work_order_status"`
	ServiceMileage             int       `json:"service_mileage" parent_entity:"trx_work_order"`
	CompanyId                  int       `json:"company" parent_entity:"trx_work_order"`
	CustomerId                 int       `json:"customer_id" parent_entity:"trx_work_order"`
	TotalAfterVAT              int       `json:"total_amount" parent_entity:"trx_work_order"`
	//BillingCustomer         int       `json:"billing_customer" gorm:"column:billable_to_id"`
}

// for sending sys_no, brand_id,model_id, vehicle chassis no to work_orderstep4

type VehicleHistoryByIdResponses struct {
	WorkOrderSystemNumber int `json:"work_order_system_number" parent_entity:"trx_work_order"`
	BrandId               int `json:"brand_id" parent_entity:"trx_work_order"`
	ModelId               int `json:"model_id" parent_entity:"trx_work_order"`
}
type VehicleHistoryGetAllResponses struct {
	WorkOrderSystemNumber   int       `json:"work_order_system_number" parent_entity:"trx_work_order"`
	WorkOrderDocumentNumber string    `json:"work_order_document_number" parent_entity:"trx_work_order" main_table:"trx_work_order"`
	WorkOrderStatusDesc     string    `json:"work_order_status_desc" parent_entity:"trx_work_order"`
	WorkOrderDate           time.Time `json:"service_request_date" parent_entity:"trx_work_order"`
	ServiceMileage          int       `json:"service_mileage" parent_entity:"trx_work_order"`
	Company                 string    `json:"company" `
	Customer                string    `json:"customer"`
	TotalAfterVAT           int       `json:"total_amount" parent_entity:"trx_work_order"`
}

type VehicleHistoryChassisResponses struct {
	VehicleChassisNo string `json:"vehicle_chassis_no"`
	VehicleEngineNo  string `json:"vehicle_engine_no"`
	Tnkb             string `json:"tnkb"`
	ModelCode        string `json:"model_code"`
	ModelName        string `json:"model_name"`
	VariantCode      string `json:"variant_code"`
	VariantName      string `json:"variant_name"`
	Status           string `json:"status"`
}

type VehicleHistoryChassisRequest struct {
	BrandId          int    `json:"brand_id"`
	ModelId          int    `json:"model_id"`
	VehicleChassisNo string `json:"vehicle_chassis_no"`
	VehicleEngineNo  string `json:"vehicle_engine_no"`
	Tnkb             string `json:"tnkb"`
	ModelCode        string `json:"model_code"`
	ModelName        string `json:"model_name"`
	VariantCode      string `json:"variant_code"`
	VariantName      string `json:"variant_name"`
	Status           string `json:"status"`
}
