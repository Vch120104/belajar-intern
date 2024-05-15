package transactionworkshopentities

const TableNameBookingEstimationServiceReminder = "trx_booking_estimation_service_reminder"

type BookingEstimationServiceReminder struct {
	BookingEstimationReminderID int    `gorm:"column:booking_estimation_reminder_id;primaryKey" json:"booking_estimation_reminder_id"`
	BookingLineNumber           int    `gorm:"column:booking_line_number;size:30;default:null" json:"booking_line_number"`
	BookingSystemNumber         int    `gorm:"column:booking_system_number;size:30;default:null" json:"booking_system_number"`
	BookingServiceReminder      string `gorm:"column:booking_service_reminder;not null" json:"booking_service_reminder"`
}

func (*BookingEstimationServiceReminder) TableName() string {
	return TableNameBookingEstimationServiceReminder
}
