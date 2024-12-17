package aftersalesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type WorkOrderDetailsData struct {
	InvoiceSystemNumber                 int     `json:"invoice_system_number"`
	WorkOrderDetailId                   int     `json:"work_order_detail_id"`
	WorkOrderSystemNumber               int     `json:"work_order_system_number"`
	LineTypeId                          int     `json:"line_type_id"`
	LineTypeCode                        string  `json:"line_type_code"`
	TransactionTypeId                   int     `json:"transaction_type_id"`
	TransactionTypeCode                 string  `json:"transaction_type_code"`
	JobTypeId                           int     `json:"job_type_id"`
	JobTypeCode                         string  `json:"job_type_code"`
	WarehouseGroupId                    int     `json:"warehouse_group_id"`
	OperationItemId                     int     `json:"operation_item_id"`
	FrtQuantity                         float64 `json:"frt_quantity"`
	SupplyQuantity                      float64 `json:"supply_quantity"`
	OperationItemPrice                  float64 `json:"operation_item_price"`
	OperationItemDiscountAmount         float64 `json:"operation_item_discount_amount"`
	OperationItemDiscountRequestAmount  float64 `json:"operation_item_discount_request_amount"`
	OperationItemDiscountPercent        float64 `json:"operation_item_discount_percent"`
	OperationItemDiscountRequestPercent float64 `json:"operation_item_discount_request_percent"`
	WorkOrderStatusId                   int     `json:"work_order_status_id"`
	TotalCostOfGoodsSoldAmount          float64 `json:"total_cost_of_goods_sold"`
	AtpmClaimNumber                     string  `json:"atpm_claim_number"`
}
type WorkOrderDetails struct {
	Data  []WorkOrderDetailsData `json:"data"`
	Limit int                    `json:"limit"`
}

type WorkOrderResponse struct {
	WorkOrderSystemNumber   int              `json:"work_order_system_number"`
	WorkOrderDocumentNumber string           `json:"work_order_document_number"`
	WorkOrderDate           string           `json:"work_order_date"`
	VehicleChassisNumber    string           `json:"vehicle_chassis_number"`
	TotalVat                float64          `json:"total_vat"`
	ProfitCenterId          int              `json:"profit_center_id"`
	WorkOrderDetails        WorkOrderDetails `json:"work_order_details"`
	CarWash                 bool             `json:"car_wash"`
	WorkOrderStatusId       int              `json:"work_order_status_id"`
	ServiceMileage          int              `json:"service_mileage"`
	DownPaymentPayment      float64          `json:"downpayment_payment"`
	DownPaymentToInvoice    float64          `json:"down_payment_to_invoice"`
	DownPaymentAllocated    float64          `json:"downpayment_payment_allocated"`
	JournalSystemNumber     int              `json:"journal_system_number"`
	BilltoCustomerId        int              `json:"billto_customer_id"`
}

type UpdateWorkOrderRequest struct {
	WorkOrderStatusId int `json:"work_order_status_id"`
}

func GetWorkOrderById(id int) (WorkOrderResponse, *exceptions.BaseErrorResponse) {
	var workOrderReq WorkOrderResponse
	url := config.EnvConfigs.AfterSalesServiceUrl + "work-order/normal/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &workOrderReq)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve work order due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "work order service is temporarily unavailable"
		}

		return workOrderReq, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting work order by ID"),
		}
	}
	return workOrderReq, nil
}
