package transactionworkshopentities

import "time"

const TableNameBookingEstimationAllocation = "trx_booking_estimation_allocation"

type BookingEstimationAllocation struct {
	BookingSystemNumber              int                                `gorm:"column:booking_system_number;size:30;primaryKey" json:"booking_system_number"`
	DocumentStatusID                 int                                `gorm:"column:document_status_id;size:30;default:null" json:"document_status_id"`
	BatchSystemNumber                int                                `gorm:"column:batch_system_number;size:30;default:null" json:"batch_system_number"`
	CompanyID                        int                                `gorm:"column:company_id;size:30;default:null" json:"company_id"`
	PdiSystemNumber                  int                                `gorm:"column:pdi_system_number;size:30;default:null" json:"pdi_system_number"`
	BookingDocumentNumber            string                             `gorm:"column:booking_document_number;type:varchar(25);unique;not null" json:"booking_document_number"`
	BookingDate                      *time.Time                          `gorm:"column:booking_date;default:null" json:"booking_date"`
	BookingStall                     string                             `gorm:"column:booking_stall;type:varchar(30);not null" json:"booking_stall"`
	BookingReminderDate              *time.Time                          `gorm:"column:booking_reminder_date;default:null" json:"booking_reminder_date"`
	BookingServiceDate               *time.Time                          `gorm:"column:booking_service_date;default:null" json:"booking_service_date"`
	BookingServiceTime               float32                            `gorm:"column:booking_service_time;type:varchar(5);default:null" json:"booking_service_time"`
	BookingEstimationTime            float32                            `gorm:"column:booking_estimation_time;default:null" json:"booking_estimation_time"`
	BookingEstimationRequest         []BookingEstimationRequest         `gorm:"foreignKey:BookingSystemNumber;" json:"booking_estimation_request"`
	BookingEstimationServiceReminder []BookingEstimationServiceReminder `gorm:"foreignKey:BookingSystemNumber;" json:"booking_estimation_service_reminder"`
}

func (*BookingEstimationAllocation) TableName() string {
	return TableNameBookingEstimationAllocation
}
