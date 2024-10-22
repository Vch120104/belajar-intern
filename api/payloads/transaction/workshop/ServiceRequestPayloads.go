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
	ServiceRequestSystemNumber int       `json:"service_request_system_number"`
	ServiceRequestDate         time.Time `json:"service_request_date"`
	ServiceRequestBy           string    `json:"service_request_by"`
	ServiceRequestStatusId     int       `json:"service_request_status_id"`
	BrandId                    int       `json:"brand_id"`
	ModelId                    int       `json:"model_id"`
	VariantId                  int       `json:"variant_id"`
	VehicleId                  int       `json:"vehicle_id"`
	BookingSystemNumber        int       `json:"booking_system_number"`
	EstimationSystemNumber     int       `json:"estimation_system_number"`
	WorkOrderSystemNumber      int       `json:"work_order_system_number"`
	ReferenceDocSystemNumber   int       `json:"reference_doc_system_number"`
	ProfitCenterId             int       `json:"profit_center_id"`
	CompanyId                  int       `json:"company_id"`
	DealerRepresentativeId     int       `json:"dealer_representative_id"`
	ServiceTypeId              int       `json:"service_type_id"`
	ReferenceTypeId            int       `json:"reference_type_id"`
	ServiceRemark              string    `json:"service_remark"`
	ServiceCompanyId           int       `json:"service_company_id"`
	ServiceDate                time.Time `json:"service_date"`
	ReplyId                    int       `json:"reply_id"`
	ServiceProfitCenterId      int       `json:"service_profit_center_id"`
	ReferenceJobType           string    `json:"reference_job_type"`
}

type ServiceRequestSaveDataRequest struct {
	ServiceTypeId    int       `json:"service_type_id"`
	ServiceRemark    string    `json:"service_remark"`
	ServiceCompanyId int       `json:"service_company_id"`
	ServiceDate      time.Time `json:"service_date"`
}

type ServiceRequestGetallResponse struct {
	ServiceRequestSystemNumber int `json:"service_request_system_number"`
	//ServiceRequestStatusId       int                           `json:"service_request_status_id"`
	ServiceRequestStatusName     string `json:"service_request_status_name"`
	ServiceRequestDocumentNumber string `json:"service_request_document_number"`
	ServiceRequestDate           string `json:"service_request_date"`
	//BrandId                      int                           `json:"brand_id"`
	BrandName string `json:"brand_name"`
	//ModelId                      int                           `json:"model_id"`
	ModelName string `json:"model_code_description"`
	//VariantId         int    `json:"variant_id"`
	VariantName       string `json:"variant_code_description"`
	VariantColourName string `json:"colour_name"`
	VehicleId         int    `json:"vehicle_id"`
	VehicleCode       string `json:"chassis_no"`
	VehicleTnkb       string `json:"no_polisi"`
	CompanyId         int    `json:"company_id"`
	CompanyName       string `json:"company_name"`
	//DealerRepresentativeId       int                           `json:"dealer_representative_id"`
	DealerRepresentativeName string `json:"dealer_rep_code_name"`
	//ProfitCenterId               int                           `json:"profit_center_id"`
	ProfitCenterName           string `json:"profit_center_name"`
	WorkOrderSystemNumber      int    `json:"work_order_no"`
	WorkOrderDocumentNumber    string `json:"work_order_document_number"`
	BookingSystemNumber        int    `json:"booking_no"`
	BookingDocumentNumber      string `json:"booking_document_number"`
	EstimationSystemNumber     int    `json:"estimation_system_number"`
	ReferenceDocSystemNumber   int    `json:"reference_doc_system_number"`
	ReferenceDocDocumentNumber string `json:"reference_doc_document_number"`
	ReferenceTypeId            int    `json:"ref_type_id"`
	ReferenceTypeName          string `json:"reference_type_name"`
	ReferenceDocId             int    `json:"ref_doc_id"`
	ReferenceDocNumber         string `json:"ref_doc_no"`
	ReferenceDocDate           string `json:"ref_doc_date"`
	//ReplyId                      int                           `json:"reply_id"`
	ReplyBy            string                        `json:"reply_by"`
	ReplyDate          string                        `json:"reply_date"`
	ReplyRemark        string                        `json:"reply_remark"`
	ServiceCompanyId   int                           `json:"service_company_id"`
	ServiceCompanyName string                        `json:"service_company_name"`
	ServiceDate        string                        `json:"service_date"`
	ServiceRequestBy   string                        `json:"service_request_by"`
	ServiceDetails     ServiceRequestDetailsResponse `json:"service_details"`
}

type ServiceRequestResponse struct {
	ServiceRequestSystemNumber int `json:"service_request_system_number"`
	//ServiceRequestStatusId       int                           `json:"service_request_status_id"`
	ServiceRequestStatusName     string `json:"service_request_status_name"`
	ServiceRequestDocumentNumber string `json:"service_request_document_number"`
	ServiceRequestDate           string `json:"service_request_date"`
	//BrandId                      int                           `json:"brand_id"`
	BrandName string `json:"brand_name"`
	//ModelId                      int                           `json:"model_id"`
	ModelName string `json:"model_code_description"`
	//VariantId         int    `json:"variant_id"`
	VariantName       string `json:"variant_code_description"`
	VariantColourName string `json:"colour_name"`
	VehicleId         int    `json:"vehicle_id"`
	VehicleCode       string `json:"chassis_no"`
	VehicleTnkb       string `json:"no_polisi"`
	CompanyId         int    `json:"company_id"`
	CompanyName       string `json:"company_name"`
	//DealerRepresentativeId       int                           `json:"dealer_representative_id"`
	DealerRepresentativeName string `json:"dealer_rep_code_name"`
	//ProfitCenterId               int                           `json:"profit_center_id"`
	ProfitCenterName         string `json:"profit_center_name"`
	ServiceProfitCenterName  string `json:"service_profit_center_name"`
	WorkOrderSystemNumber    int    `json:"work_order_no"`
	WorkOrderDocumentNumber  string `json:"work_order_document_number"`
	BookingSystemNumber      int    `json:"booking_no"`
	EstimationSystemNumber   int    `json:"estimation_system_number"`
	ReferenceDocSystemNumber int    `json:"reference_doc_system_number"`
	ReferenceTypeId          int    `json:"ref_type_id"`
	ReferenceTypeName        string `json:"reference_type_name"`
	ReferenceDocId           int    `json:"ref_doc_id"`
	ReferenceDocNumber       string `json:"ref_doc_no"`
	ReferenceDocDate         string `json:"ref_doc_date"`
	ServiceRemark            string `json:"service_remark"`
	//ReplyId                      int                           `json:"reply_id"`
	ReplyBy            string                        `json:"reply_by"`
	ReplyDate          string                        `json:"reply_date"`
	ReplyRemark        string                        `json:"reply_remark"`
	ServiceCompanyId   int                           `json:"service_company_id"`
	ServiceCompanyName string                        `json:"service_company_name"`
	ServiceDate        string                        `json:"service_date"`
	ServiceRequestBy   string                        `json:"service_request_by"`
	ServiceDetails     ServiceRequestDetailsResponse `json:"service_details"`
}

type SubmitServiceRequestResponse struct {
	DocumentNumber             string `json:"service_request_document_number"`
	ServiceRequestSystemNumber int    `json:"service_request_system_number"`
}

type WorkOrderRequestResponse struct {
	WorkOrderDocumentNumber string `json:"work_order_document_number"`
	WorkOrderSystemNumber   int    `json:"work_order_system_number"`
}

type ServiceDetailSaveRequest struct {
	ServiceRequestSystemNumber int     `json:"service_request_system_number"`
	LineTypeId                 int     `json:"line_type_id"`
	OperationItemId            int     `json:"operation_item_id"`
	FrtQuantity                float64 `json:"frt_quantity"`
}

type ServiceDetailUpdateRequest struct {
	ServiceRequestSystemNumber int     `json:"service_request_system_number"`
	FrtQuantity                float64 `json:"frt_quantity"`
}

type ServiceRequestDetail struct {
	ServiceRequestDetailId     int     `json:"service_request_detail_id" parent_entity:"trx_service_request_detail" main_table:"trx_service_request_detail"`
	ServiceRequestLineNumber   int     `json:"service_request_line_number" parent_entity:"trx_service_request_detail" `
	ServiceRequestSystemNumber int     `json:"service_request_system_number" parent_entity:"trx_service_request_detail" `
	LineTypeId                 int     `json:"line_type_id" parent_entity:"trx_service_request_detail" `
	OperationItemId            int     `json:"operation_item_id" parent_entity:"trx_service_request_detail" `
	FrtQuantity                float64 `json:"frt_quantity" parent_entity:"trx_service_request_detail" `
}

type ServiceRequestDetailsResponse struct {
	Page       int                     `json:"page"`
	Limit      int                     `json:"limit"`
	TotalPages int                     `json:"total_pages"`
	TotalRows  int                     `json:"total_rows"`
	Data       []ServiceDetailResponse `json:"data"`
}

type ServiceDetailResponse struct {
	ServiceRequestDetailId     int     `json:"service_request_detail_id"`
	ServiceRequestSystemNumber int     `json:"service_request_system_number"`
	LineTypeId                 int     `json:"line_type_id"`
	LineTypeCode               string  `json:"line_type_code"`
	OperationItemId            int     `json:"operation_item_id"`
	OperationItemCode          string  `json:"code"`
	OperationItemName          string  `json:"description"`
	UomName                    string  `json:"uom"`
	FrtQuantity                float64 `json:"qty"`
	ReferenceDocSystemNumber   int     `json:"reference_doc_system_number"`
	ReferenceDocCode           string  `json:"reference_doc_code"`
	ReferenceDocNumber         string  `json:"reference_doc_name"`
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

type ServiceRequestBookingEstimation struct {
	ProfitCenterId               int    `json:"profit_center_id"`
	CompanyId                    int    `json:"company_id"`
	VehicleId                    int    `json:"vehicle_id"`
	ServiceRequestDocumentNumber string `json:"service_request_document_number"`
	ContractServiceSystemNumber  int    `json:"contract_service_system_number"`
}

type ServiceRequestDetailBookingPayloads struct {
	ServiceRequestDetailId     int     `json:"service_request_detail_id"`
	ServiceRequestSystemNumber int     `json:"service_request_system_number"`
	ReferenceDocumentNumber    string  `json:"reference_document_number"`
	LineTypeId                 int     `json:"line_type_id"`
	OperationItemId            int     `json:"operation_item_id"`
	FrtQuantity                float64 `json:"frt_quantity"`
}
type ProfitCenter struct {
	ProfitCenterId   int    `json:"profit_center_id"`
	ProfitCenterCode string `json:"profit_center_code"`
	ProfitCenterName string `json:"profit_center_name"`
}

type DealerRepresentative struct {
	DealerRepresentativeId   int    `json:"dealer_representative_id"`
	DealerRepresentativeName string `json:"dealer_representative_name"`
}

type ReferenceType struct {
	ReferenceTypeId   int    `json:"service_request_reference_type_id"`
	ReferenceTypeCode string `json:"service_request_reference_type_code"`
	ReferenceTypeName string `json:"service_request_reference_type_description"`
}

type ReferenceDoc struct {
	ReferenceDocSystemNumber int    `json:"reference_doc_system_number"`
	ReferenceDocNumber       string `json:"reference_doc_number"`
	ReferenceDocCode         string `json:"reference_doc_code"`
	ReferenceDocDate         string `json:"reference_doc_date"`
}

type ServiceType struct {
	ServiceTypeId   int    `json:"service_profit_center_id"`
	ServiceTypeCode string `json:"service_profit_center_code"`
	ServiceTypeName string `json:"service_profit_center_description"`
}
