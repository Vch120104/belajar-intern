package transactionworkshoppayloads

import "time"

type WorkOrderDetailBypassRequest struct {
	WorkOrderDetailId     int     `json:"work_order_detail_id" parent_entity:"trx_work_order_detail" main_table:"trx_work_order_detail"`
	WorkOrderSystemNumber int     `json:"work_order_system_number" parent_entity:"trx_work_order_detail"`
	LineTypeId            int     `json:"line_type_id" parent_entity:"trx_work_order_detail"`
	TransactionTypeId     int     `json:"transaction_type_id" parent_entity:"trx_work_order_detail" `
	JobTypeId             int     `json:"job_type_id" parent_entity:"trx_work_order_detail"`
	FrtQuantity           float64 `json:"frt_quantity" parent_entity:"trx_work_order_detail"`
	SupplyQuantity        float64 `json:"supply_quantity" parent_entity:"trx_work_order_detail"`
	PriceListId           int     `json:"price_list_id" parent_entity:"trx_work_order_detail"`
	WarehouseId           int     `json:"warehouse_id" parent_entity:"trx_work_order_detail"`
	ItemId                int     `json:"item_id" parent_entity:"trx_work_order_detail"`
	ProposedPrice         float64 `json:"operation_item_discount_request_amount" parent_entity:"trx_work_order_detail"`
	OperationItemPrice    float64 `json:"operation_item_price" parent_entity:"trx_work_order_detail"`
}

type WorkOrderDetailBypassResponse struct {
	WorkOrderDetailId                  int     `json:"work_order_detail_id"`
	WorkOrderSystemNumber              int     `json:"work_order_system_number"`
	WorkOrderDocumentNumber            string  `json:"work_order_document_number"`
	LineTypeId                         int     `json:"line_type_id"`
	LineTypeName                       string  `json:"line_type_name"`
	TransactionTypeId                  int     `json:"transaction_type_id"`
	JobTypeId                          int     `json:"job_type_id"`
	WarehouseId                        int     `json:"warehouse_id"`
	ItemId                             int     `json:"item_id"`
	ItemCode                           string  `json:"item_code"`
	ItemName                           string  `json:"item_name"`
	FrtQuantity                        float64 `json:"frt_quantity"`
	SupplyQuantity                     float64 `json:"supply_quantity"`
	OperationItemPrice                 float64 `json:"operation_item_price"`
	OperationItemDiscountAmount        float64 `json:"operation_item_discount_amount"`
	OperationItemDiscountRequestAmount float64 `json:"operation_item_discount_request_amount"`
}

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
