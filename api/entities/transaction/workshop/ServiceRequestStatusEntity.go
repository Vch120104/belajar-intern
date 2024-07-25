package transactionworkshopentities

var CreateServiceRequestMasterStatusTable = "mtr_service_request_status"

type ServiceRequestMasterStatus struct {
	ServiceRequestStatusId          int    `gorm:"column:service_request_status_id;size:30;not null;primaryKey" json:"service_request_status_id"`
	ServiceRequestStatusCode        string `gorm:"column:service_request_status_code;size:30;" json:"service_request_status_code"`
	ServiceRequestStatusDescription string `gorm:"column:service_request_status_description;size:30;" json:"service_request_status_description"`
}

func (*ServiceRequestMasterStatus) TableName() string {
	return CreateServiceRequestMasterStatusTable
}
