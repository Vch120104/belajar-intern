package transactionworkshopentities

import "math"

const TableNameServiceRequestDetail = "trx_service_request_detail"

type ServiceRequestDetail struct {
	ServiceRequestDetailId     int     `gorm:"column:service_request_detail_id;size:30;primary_key;" json:"service_request_detail_id"`
	ServiceRequestLineNumber   int     `gorm:"column:service_request_line_number;size:30;" json:"service_request_line_number"`
	ServiceRequestSystemNumber int     `gorm:"column:service_request_system_number;size:30;" json:"service_request_system_number"`
	ReferenceDocSystemNumber   int     `gorm:"column:reference_doc_system_number;size:30;" json:"reference_doc_system_number"`
	ReferenceDocId             int     `gorm:"column:reference_doc_id;size:30;" json:"reference_doc_id"`
	LineTypeId                 int     `gorm:"column:line_type_id;size:30;" json:"line_type_id"`
	OperationItemId            int     `gorm:"column:operation_item_id;size:30;" json:"operation_item_id"`
	FrtQuantity                float64 `gorm:"column:frt_quantity;" json:"frt_quantity"`
}

// SetFrtQuantity sets the FrtQuantity with scale 2
func (s *ServiceRequestDetail) SetFrtQuantity(value float64) {
	s.FrtQuantity = math.Round(value*100) / 100
}

// GetFrtQuantity gets the FrtQuantity
func (s *ServiceRequestDetail) GetFrtQuantity() float64 {
	return s.FrtQuantity
}

func (*ServiceRequestDetail) TableName() string {
	return TableNameServiceRequestDetail
}
