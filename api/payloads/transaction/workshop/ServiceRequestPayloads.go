package transactionworkshoppayloads

import "time"

type ServiceRequestNew struct {
	ServiceRequestSystemNumber   int       `json:"service_request_system_number" parent_entity:"trx_service_request" main_table:"trx_service_request"`
	ServiceRequestDocumentNumber string    `json:"service_request_document_number" parent_entity:"trx_service_request" `
	ServiceRequestDate           time.Time `json:"service_request_date" parent_entity:"trx_service_request" `
	ServiceRequestBy             string    `json:"service_request_by" parent_entity:"trx_service_request" `
	ServiceRequestStatusId       int       `json:"service_request_status_id" parent_entity:"trx_service_request" `
	BrandId                      int       `json:"brand_id" parent_entity:"trx_service_request" `
	ModelId                      int       `json:"model_id" parent_entity:"trx_service_request" `
	VariantId                    int       `json:"variant_id" parent_entity:"trx_service_request" `
	VehicleId                    int       `json:"vehicle_id" parent_entity:"trx_service_request" `
	BookingSystemNumber          int       `json:"booking_system_number" parent_entity:"trx_service_request" `
	EstimationSystemNumber       int       `json:"estimation_system_number" parent_entity:"trx_service_request" `
	WorkOrderSystemNumber        int       `json:"work_order_system_number" parent_entity:"trx_service_request" `
	ReferenceDocSystemNumber     int       `json:"reference_doc_system_number" parent_entity:"trx_service_request" `
	ProfitCenterId               int       `json:"profit_center_id" parent_entity:"trx_service_request" `
	CompanyId                    int       `json:"company_id" parent_entity:"trx_service_request" `
	DealerRepresentativeId       int       `json:"dealer_representative_id" parent_entity:"trx_service_request" `
	ServiceTypeId                int       `json:"service_type_id" parent_entity:"trx_service_request" `
	ReferenceTypeId              int       `json:"reference_type_id" parent_entity:"trx_service_request" `
	ServiceRemark                string    `json:"service_remark" parent_entity:"trx_service_request" `
	ServiceCompanyId             int       `json:"service_company_id" parent_entity:"trx_service_request" `
	ServiceDate                  time.Time `json:"service_date" parent_entity:"trx_service_request" `
	ReplyId                      int       `json:"reply_id" parent_entity:"trx_service_request" `
}

type ServiceRequestSaveRequest struct {
	ServiceRequestSystemNumber   int       `json:"service_request_system_number"`
	ServiceRequestDocumentNumber string    `json:"service_request_document_number"`
	ServiceRequestDate           time.Time `json:"service_request_date"`
	ServiceRequestBy             string    `json:"service_request_by"`
	ServiceRequestStatusId       int       `json:"service_request_status_id"`
	BrandId                      int       `json:"brand_id"`
	ModelId                      int       `json:"model_id"`
	VariantId                    int       `json:"variant_id"`
	VehicleId                    int       `json:"vehicle_id"`
	BookingSystemNumber          int       `json:"booking_system_number"`
	EstimationSystemNumber       int       `json:"estimation_system_number"`
	WorkOrderSystemNumber        int       `json:"work_order_system_number"`
	ReferenceDocSystemNumber     int       `json:"reference_doc_system_number"`
	ProfitCenterId               int       `json:"profit_center_id"`
	CompanyId                    int       `json:"company_id"`
	DealerRepresentativeId       int       `json:"dealer_representative_id"`
	ServiceTypeId                int       `json:"service_type_id"`
	ReferenceTypeId              int       `json:"reference_type_id"`
	ServiceRemark                string    `json:"service_remark"`
	ServiceCompanyId             int       `json:"service_company_id"`
	ServiceDate                  time.Time `json:"service_date"`
	ReplyId                      int       `json:"reply_id"`
	ServiceProfitCenterId        int       `json:"service_profit_center_id"`
	ReferenceJobType             string    `json:"reference_job_type"`
}

type ServiceRequestResponse struct {
	ServiceRequestSystemNumber   int    `json:"service_request_system_number"`
	ServiceRequestDocumentNumber string `json:"service_request_document_number"`
	ServiceRequestDate           string `json:"service_request_date"`
	ServiceRequestBy             string `json:"service_request_by"`
	ServiceRequestStatusId       int    `json:"service_request_status_id"`
	ServiceRequestStatusName     string `json:"service_request_status_description"`
	BrandId                      int    `json:"brand_id"`
	BrandName                    string `json:"brand_name"`
	ModelId                      int    `json:"model_id"`
	ModelName                    string `json:"model_name"`
	VariantId                    int    `json:"variant_id"`
	VariantName                  string `json:"variant_name"`
	VariantColourName            string `json:"variant_colour_name"`
	VehicleId                    int    `json:"vehicle_id"`
	VehicleCode                  string `json:"vehicle_chassis_number"`
	VehicleTnkb                  string `json:"vehicle_registration_certificate_tnkb"`
	BookingSystemNumber          int    `json:"booking_system_number"`
	EstimationSystemNumber       int    `json:"estimation_system_number"`
	WorkOrderSystemNumber        int    `json:"work_order_system_number"`
	WorkOrderDocumentNumber      string `json:"work_order_document_number"`
	ReferenceDocSystemNumber     int    `json:"reference_doc_system_number"`
	ProfitCenterId               int    `json:"profit_center_id"`
	CompanyId                    int    `json:"company_id"`
	CompanyName                  string `json:"company_name"`
	DealerRepresentativeId       int    `json:"dealer_representative_id"`
	ServiceTypeId                int    `json:"service_type_id"`
	ReferenceTypeId              int    `json:"reference_type_id"`
	ServiceRemark                string `json:"service_remark"`
	ServiceCompanyId             int    `json:"service_company_id"`
	ServiceCompanyName           string `json:"service_company_name"`
	ServiceDate                  string `json:"service_date"`
	ReplyId                      int    `json:"reply_id"`
	ReplyDate                    string `json:"reply_date"`
	ReplyBy                      string `json:"reply_by"`
	ReplyRemark                  string `json:"reply_remark"`
}

// type ServiceRequestDetailsResponse struct {
// 	// Existing fields
// 	ServiceDetails []ServiceDetailResponse `json:"service_details"` // New field for details
// }

type SubmitServiceRequestResponse struct {
	DocumentNumber             string `json:"service_request_document_number"`
	ServiceRequestSystemNumber int    `json:"service_request_system_number"`
}

type WorkOrderRequestResponse struct {
	WorkOrderDocumentNumber string `json:"work_order_document_number"`
	WorkOrderSystemNumber   int    `json:"work_order_system_number"`
}

type ServiceDetailSaveRequest struct {
	ServiceRequestDetailId     int     `json:"service_request_detail_id"`
	ServiceRequestId           int     `json:"service_request_id"`
	ServiceRequestSystemNumber int     `json:"service_request_system_number"`
	LineTypeId                 int     `json:"line_type_id"`
	OperationItemId            int     `json:"operation_item_id"`
	ReferenceDocSystemNumber   int     `json:"reference_doc_system_number"`
	ReferenceDocId             int     `json:"reference_doc_id"`
	FrtQuantity                float64 `json:"frt_quantity"`
}

type ServiceDetailResponse struct {
	ServiceRequestDetailId     int     `json:"service_request_detail_id"`
	ServiceRequestId           int     `json:"service_request_id"`
	ServiceRequestSystemNumber int     `json:"service_request_system_number"`
	LineTypeId                 int     `json:"line_type_id"`
	OperationItemId            int     `json:"operation_item_id"`
	FrtQuantity                float64 `json:"frt_quantity"`
	ReferenceDocSystemNumber   int     `json:"reference_doc_system_number"`
	ReferenceDocId             int     `json:"reference_doc_id"`
}

type ServiceRequestDetail struct {
	ServiceRequestDetailId     int     `json:"service_request_detail_id" parent_entity:"trx_service_request_detail" main_table:"trx_service_request_detail"`
	ServiceRequestId           int     `json:"service_request_id" parent_entity:"trx_service_request_detail" `
	ServiceRequestSystemNumber int     `json:"service_request_system_number" parent_entity:"trx_service_request_detail" `
	LineTypeId                 int     `json:"line_type_id" parent_entity:"trx_service_request_detail" `
	OperationItemId            int     `json:"operation_item_id" parent_entity:"trx_service_request_detail" `
	FrtQuantity                float64 `json:"frt_quantity" parent_entity:"trx_service_request_detail" `
}

type ServiceRequestDetailResponse struct {
	ServiceRequestDetailId     int     `json:"service_request_detail_id"`
	ServiceRequestId           int     `json:"service_request_id"`
	ServiceRequestSystemNumber int     `json:"service_request_system_number"`
	LineTypeId                 int     `json:"line_type_id"`
	OperationItemId            int     `json:"operation_item_id"`
	FrtQuantity                float64 `json:"frt_quantity"`
	ReferenceDocSystemNumber   int     `json:"reference_doc_system_number"`
	ReferenceDocId             int     `json:"reference_doc_id"`
}
type ItemServiceRequestDetail struct {
	ItemId   int    `json:"item_id"`
	ItemName string `json:"item_name"`
	ItemCode string `json:"item_code"`
	UomId    int    `json:"unit_of_measurement_type_id"`
}

type UomItemServiceRequestDetail struct {
	UomId   int    `json:"uom_id"`
	UomName string `json:"uom_description"`
}

type CompanyResponse struct {
	CompanyId   int    `json:"company_id"`
	CompanyName string `json:"company_name"`
	BizCategory string `json:"biz_category"`
}

type ServiceRequestStatusResponse struct {
	ServiceRequestStatusId   int    `json:"service_request_status_id"`
	ServiceRequestStatusCode string `json:"service_request_status_code"`
	ServiceRequestStatusName string `json:"service_request_status_description"`
}
