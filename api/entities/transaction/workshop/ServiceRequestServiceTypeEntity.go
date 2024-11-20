package transactionworkshopentities

var CreateServiceRequestServiceTypeTable = "mtr_service_request_service_type"

type ServiceRequestMasterServiceType struct {
	ServiceRequestServiceTypeId          int    `gorm:"column:service_request_service_type_id;size:30;not null;primaryKey" json:"service_request_service_type_id"`
	ServiceRequestServiceTypeCode        string `gorm:"column:service_request_service_type_code;size:30;" json:"service_request_service_type_code"`
	ServiceRequestServiceTypeDescription string `gorm:"column:service_request_service_type_description;size:30;" json:"service_request_service_type_description"`
}

func (*ServiceRequestMasterServiceType) TableName() string {
	return CreateServiceRequestServiceTypeTable
}
