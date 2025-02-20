package transactionworkshopentities

const TableNameBookingEstimationRequest = "trx_booking_estimation_request"

type BookingEstimationRequest struct {
	BookingEstimationRequestID int    `gorm:"column:booking_estimation_request_id;size:30;primaryKey" json:"booking_estimation_request_id"`
	BookingSystemNumber        int    `gorm:"column:booking_system_number;size:30;default:null" json:"booking_system_number"`
	BookingDocumentNumber      string `gorm:"column:booking_document_number;" json:"booking_document_number"`
	BookingServiceRequest      string `gorm:"column:booking_service_request;not null" json:"booking_service_request"`
	BookingLine                int    `gorm:"column:booking_line;size:30;" json:"booking_line"`
	IsActive                   bool   `gorm:"column:is_active;default:true" json:"is_active"`
}

func (*BookingEstimationRequest) TableName() string {
	return TableNameBookingEstimationRequest
}
