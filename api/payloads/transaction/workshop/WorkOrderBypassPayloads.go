package transactionworkshoppayloads

import "time"

type WorkOrderBypassRequest struct {
	WorkOrderDetailId     int `json:"work_order_detail_id"`
	WorkOrderSystemNumber int `json:"work_order_system_number"`
	LineTypeId            int `json:"line_type_id"`
	ItemId                int `json:"item_id"`
	OperationId           int `json:"operation_id"`
	ServiceAdvisorId      int `json:"service_advisor_id"`
	TechnicianId          int `json:"technician_id"`
	ForemanWSId           int `json:"foreman_ws_id"`
	ForemanBSId           int `json:"foreman_bs_id"`
}

type WorkOrderBypassResponse struct {
	WorkOrderDocumentNumber string `json:"work_order_document_number"`
	WorkOrderSystemNumber   int    `json:"work_order_system_number"`
	LineTypeId              int    `json:"line_type_id"`
	LineTypeName            string `json:"line_type_name"`
	ItemId                  int    `json:"item_id"`
	ItemName                string `json:"item_name"`
	OperationId             int    `json:"operation_id"`
	OperationName           string `json:"operation_name"`
	ServiceAdvisorId        int    `json:"service_advisor_id"`
	ServiceAdvisorName      string `json:"service_advisor_name"`
	TechnicianId            int    `json:"technician_id"`
	TechnicianName          string `json:"technician_name"`
	ForemanWSId             int    `json:"foreman_ws_id"`
	ForemanWSName           string `json:"foreman_ws_name"`
	ForemanBSId             int    `json:"foreman_bs_id"`
	ForemanBSName           string `json:"foreman_bs_name"`
}

type WorkOrderBypassRequestDetail struct {
	WorkOrderSystemNumber           int       `json:"work_order_system_number"`
	WorkOrderDocumentNumber         string    `json:"work_order_document_number"`
	WorkOrderQualityControlStatusID int       `json:"work_order_quality_control_status_id"`
	WorkOrderStartDateTime          time.Time `json:"work_order_start_date_time"`
	WorkOrderEndDateTime            time.Time `json:"work_order_end_date_time"`
	WorkOrderActualTime             float32   `json:"work_order_actual_time"`
}

type WorkOrderBypassResponseDetail struct {
	WorkOrderSystemNumber           int       `json:"work_order_system_number"`
	WorkOrderDocumentNumber         string    `json:"work_order_document_number"`
	WorkOrderQualityControlStatusID int       `json:"work_order_quality_control_status_id"`
	WorkOrderStartDateTime          time.Time `json:"work_order_start_date_time"`
	WorkOrderEndDateTime            time.Time `json:"work_order_end_date_time"`
	WorkOrderActualTime             float32   `json:"work_order_actual_time"`
}
