package transactionunitentities

import "time"

var CreatePDIRequestDetailTable = "dms_microservices_sales_dev.dbo.trx_pdi_request_detail"

type PdiRequestDetail struct {
	PdiRequestDetailSystemNumber int       `gorm:"column:pdi_request_detail_system_number;primaryKey;autoIncrement" json:"pdi_request_detail_system_number"` // Primary key dan auto increment
	PdiRequestDetailLineNumber   int       `gorm:"column:pdi_request_detail_line_number;not null" json:"pdi_request_detail_line_number"`                     // Kolom wajib
	PdiRequestSystemNumber       *int      `gorm:"column:pdi_request_system_number" json:"pdi_request_system_number"`                                        // Nullable
	OperationNumberId            int       `gorm:"column:operation_number_id;not null" json:"operation_number_id"`                                           // Kolom wajib
	VehicleId                    *int      `gorm:"column:vehicle_id" json:"vehicle_id"`                                                                      // Nullable
	EstimatedDelivery            time.Time `gorm:"column:estimated_delivery;not null" json:"estimated_delivery"`                                             // Kolom wajib
	LineRemark                   *string   `gorm:"column:pdi_request_detail_line_remark;size:256" json:"line_remark"`                                        // Nullable
	Frt                          float64   `gorm:"column:frt;not null" json:"frt"`                                                                           // Kolom wajib
	ServiceDate                  time.Time `gorm:"column:service_date;not null" json:"service_date"`                                                         // Kolom wajib
	ServiceTime                  time.Time `gorm:"column:service_time;not null" json:"service_time"`                                                         // Kolom wajib
	PdiRequestDetailLineStatusId int       `gorm:"column:pdi_request_detail_line_status_id;not null" json:"pdi_request_detail_line_status_id"`               // Kolom wajib
	BookingSystemNumber          *int      `gorm:"column:booking_system_number" json:"booking_system_number"`                                                // Nullable
	WorkOrderSystemNumber        *int      `gorm:"column:work_order_system_number" json:"work_order_system_number"`                                          // Nullable
	InvoicePayableSystemNumber   *int      `gorm:"column:invoice_payable_system_number" json:"invoice_payable_system_number"`                                // Nullable
}
