package transactionunitentities

import "time"

var CreatePDIRequestDetailTable = "trx_pdi_request_detail"

type PdiRequestDetail struct {
	PdiRequestDetailSystemNumber int `gorm:"column:pdi_request_detail_system_number;size:30" json:"pdi_request_detail_system_number"`
	PdiRequestDetailLineNumber   int `gorm:"column:pdi_request_detail_line_number;size:30" json:"pdi_request_detail_line_number"`
	PdiRequestSystemNumber       int `gorm:"column:pdi_request_system_number;size:30" json:"pdi_request_system_number"`
	OperationNumberId            int `gorm:"column:operation_number_id;size:30" json:"operation_number_id"`
	VehicleId                    int `gorm:"column:vehicle_id;size:30;not null" json:"vehicle_id"`
	EstimatedDelivery            time.Time`gorm:"column:estimated_delivery;not null" json:"estimated_delivery"`
	LineRemark string `gorm:"column:line_remark;size:256;not null" json:"line_remark"`
	Frt float64 `gorm:"column:frt;not null" json:"frt"`
	ServiceDate time.Time `gorm:"column:service_date;not null" json:"service_date"`
	ServiceTime float64 `gorm:"column:service_time;not null" json:"service_time"`
	PdiRequestDetailLineStatusId int `gorm:"column:pdi_request_detail_line_status_id;size:30"`
	BookingSystemNumber int `gorm:"column:booking_system_number;size:30" json:"bookin_system_number"`
	WorkOrderSystemNumber int `gorm:"column:workorder_system_number;size:30" json:"workorder_system_number"`
	InvoicePayableSystemNumber int `gorm:"column:invoice_payable_system_number;size:30" json:"invoice_payable_system_number"`
}