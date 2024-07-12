package transactionworkshopentities

const TableNameBookingEstimationRequest = "trx_booking_estimation_request"

type BookingEstimationRequest struct {
	BookingEstimationRequestID   int    `gorm:"column:booking_estimation_request_id;primaryKey" json:"booking_estimation_request_id"`
	BookingEstimationRequestCode int    `gorm:"column:booking_estimation_request_code;size:30;default:null" json:"booking_estimation_request_code"`
	BookingSystemNumber          int    `gorm:"column:booking_system_number;size:30;default:null" json:"booking_system_number"`
	BookingServiceRequest        string `gorm:"column:booking_service_request;not null" json:"booking_service_request"`
}

func (*BookingEstimationRequest) TableName() string {
	return TableNameBookingEstimationRequest
}
