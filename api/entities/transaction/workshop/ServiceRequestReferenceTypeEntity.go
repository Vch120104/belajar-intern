package transactionworkshopentities

var CreateServiceRequestReferenceTypeTable = "mtr_service_request_reference_type"

type ServiceRequestMasterReferenceType struct {
	ServiceRequestReferenceTypeId          int    `gorm:"column:service_request_reference_type_id;size:30;not null;primaryKey" json:"service_request_reference_type_id"`
	ServiceRequestReferenceTypeCode        string `gorm:"column:service_request_reference_type_code;size:30;" json:"service_request_reference_type_code"`
	ServiceRequestReferenceTypeDescription string `gorm:"column:service_request_reference_type_description;size:30;" json:"service_request_reference_type_description"`
}

func (*ServiceRequestMasterReferenceType) TableName() string {
	return CreateServiceRequestReferenceTypeTable
}
