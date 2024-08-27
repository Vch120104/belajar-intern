package transactionworkshoppayloads

import "time"

type ServiceWorkshopRequest struct {
	ServiceLogSystemNumber           int       `json:"service_log_system_number" parent_entity:"trx_service_log" main_table:"trx_service_log"`
	WorkOrderSystemNumber            int       `json:"work_order_system_number" parent_entity:"trx_service_log"`
	TechnicianAllocationLine         int       `json:"technician_allocation_line" parent_entity:"trx_service_log"`
	TechnicianAllocationSystemNumber int       `json:"technician_allocation_system_number" parent_entity:"trx_service_log"`
	TechnicianId                     int       `json:"technician_id" parent_entity:"trx_service_log"`
	ServiceStatusId                  int       `json:"service_status_id" parent_entity:"trx_service_log"`
	ActualTime                       float64   `json:"actual_time" parent_entity:"trx_service_log"`
	PendingTime                      float64   `json:"pending_time" parent_entity:"trx_service_log"`
	WorkOrderOperationId             int       `json:"work_order_operation_id" parent_entity:"trx_service_log"`
	ShiftScheduleId                  int       `json:"shift_schedule_id" parent_entity:"trx_service_log"`
	StartDatetime                    time.Time `json:"start_datetime" parent_entity:"trx_service_log"`
	EndDatetime                      time.Time `json:"end_datetime" parent_entity:"trx_service_log"`
	EstimatedPendingTime             float64   `json:"estimated_pending_time" parent_entity:"trx_service_log"`
	Remark                           string    `json:"remark" parent_entity:"trx_service_log"`
	ActualStartTime                  float64   `json:"actual_start_time" parent_entity:"trx_service_log"`
}

type ServiceWorkshopResponse struct {
	TechnicianAllocationSystemNumber string    `json:"technician_allocation_system_number"`
	StartDatetime                    time.Time `json:"start_datetime"`
	OperationItemCode                string    `json:"operation_code"`
	Frt                              float64   `json:"frt"`
	ServiceStatusId                  int       `json:"service_status_id"`
	ServiceStatusDescription         string    `json:"service_status_description"`
	ServActualTime                   float64   `json:"serv_actual_time"`
	ServPendingTime                  float64   `json:"serv_pending_time"`
	ServProgressTime                 float64   `json:"serv_progress_time"`
	TechAllocStartDate               string    `json:"tech_alloc_start_date"`
	TechAllocStartTime               string    `json:"tech_alloc_start_time"`
	TechAllocEndDate                 string    `json:"tech_alloc_end_date"`
	TechAllocEndTime                 string    `json:"tech_alloc_end_time"`
}

// TimeReference represents the structure of the data returned by the API
type TimeReference struct {
	TimeDiff int `json:"time_different"` // `json:"time_different"`
}

type ServiceWorkshopDetailResponse struct {
	ServiceTypeName         string                         `json:"service_type_name"`
	ForemanName             string                         `json:"foreman_name"`
	TechnicianId            int                            `json:"technician_id"`
	TechnicianName          string                         `json:"technician_name"`
	WorkOrderSystemNumber   int                            `json:"work_order_system_number"`
	WorkOrderDocumentNumber string                         `json:"work_order_document_number"`
	WorkOrderDate           string                         `json:"work_order_date"`
	ModelName               string                         `json:"model_name"`
	VariantName             string                         `json:"variant_name"`
	VehicleCode             string                         `json:"vehicle_chassis_number"`
	VehicleTnkb             string                         `json:"vehicle_tnkb"`
	ServiceDetails          ServiceWorkshopDetailsResponse `json:"service_details"`
}

type ServiceWorkshopDetailsResponse struct {
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"total_pages"`
	TotalRows  int                       `json:"total_rows"`
	Data       []ServiceWorkshopResponse `json:"data"`
}

type ServiceWorkshopWoResponse struct {
	WorkOrderSystemNumber   int    `json:"work_order_system_number"`
	WorkOrderDocumentNumber string `json:"work_order_document_number"`
	WorkOrderDate           string `json:"work_order_date"`
	ModelName               string `json:"model_name"`
	VariantName             string `json:"variant_name"`
	VehicleCode             string `json:"vehicle_chassis_number"`
	VehicleTnkb             string `json:"vehicle_tnkb"`
}

type ServiceStatusResponse struct {
	ServiceStatusId   int    `json:"service_status_id"`
	ServiceStatusCode string `json:"service_status_code"`
	ServiceStatusName string `json:"service_status_description"`
}
