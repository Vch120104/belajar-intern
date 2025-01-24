package transactionworkshopentities

import "time"

const TableNameServiceRequest = "trx_service_request"

type ServiceRequest struct {
	ServiceRequestSystemNumber   int                    `gorm:"column:service_request_system_number;size:30;primary_key;" json:"service_request_system_number"`
	ServiceRequestDocumentNumber string                 `gorm:"column:service_request_document_number;" json:"service_request_document_number"`
	ServiceRequestDate           time.Time              `gorm:"column:service_request_date;size:30;" json:"service_request_date"`
	ServiceRequestBy             string                 `gorm:"column:service_request_by;size:30;" json:"service_request_by"`
	ServiceRequestStatusId       int                    `gorm:"column:service_request_status_id;size:30;" json:"service_request_status_id"`
	BrandId                      int                    `gorm:"column:brand_id;size:30;" json:"brand_id"`
	ModelId                      int                    `gorm:"column:model_id;size:30;" json:"model_id"`
	VariantId                    int                    `gorm:"column:variant_id;size:30;" json:"variant_id"`
	ColourId                     int                    `gorm:"column:colour_id;size:30;" json:"colour_id"`
	VehicleId                    int                    `gorm:"column:vehicle_id;size:30;" json:"vehicle_id"`
	BookingSystemNumber          int                    `gorm:"column:booking_system_number;size:30;" json:"booking_system_number"`
	EstimationSystemNumber       int                    `gorm:"column:estimation_system_number;size:30;" json:"estimation_system_number"`
	WorkOrderSystemNumber        int                    `gorm:"column:work_order_system_number;size:30;" json:"work_order_system_number"`
	ReferenceSystemNumber        int                    `gorm:"column:reference_system_number;size:30;" json:"reference_system_number"`
	ProfitCenterId               int                    `gorm:"column:profit_center_id;size:30;" json:"profit_center_id"`
	CompanyId                    int                    `gorm:"column:company_id;size:30;" json:"company_id"`
	DealerRepresentativeId       int                    `gorm:"column:dealer_representative_id;size:30;" json:"dealer_representative_id"`
	ServiceTypeId                int                    `gorm:"column:service_type_id;size:30;" json:"service_type_id"`
	ReferenceTypeId              int                    `gorm:"column:reference_type_id;size:30;" json:"reference_type_id"`
	ReferenceJobType             string                 `gorm:"column:reference_job_type;size:30;" json:"reference_job_type"`
	ServiceRemark                string                 `gorm:"column:service_remark;size:30;" json:"service_remark"`
	ServiceCompanyId             int                    `gorm:"column:service_company_id;size:30;" json:"service_company_id"`
	ServiceDate                  time.Time              `gorm:"column:service_date;size:30;" json:"service_date"`
	ReplyId                      int                    `gorm:"column:reply_id;size:30;" json:"reply_id"`
	ReplyDate                    time.Time              `gorm:"column:reply_date;size:30;" json:"reply_date"`
	ReplyBy                      string                 `gorm:"column:reply_by;size:30;" json:"reply_by"`
	ReplyRemark                  string                 `gorm:"column:reply_remark;size:30;" json:"reply_remark"`
	ServiceProfitCenterId        int                    `gorm:"column:service_profit_center_id;size:30;" json:"service_profit_center_id"`
	ServiceRequestDetail         []ServiceRequestDetail `gorm:"foreignKey:ServiceRequestSystemNumber;references:ServiceRequestSystemNumber" json:"service_request_detail"`
}

func (*ServiceRequest) TableName() string {
	return TableNameServiceRequest
}
