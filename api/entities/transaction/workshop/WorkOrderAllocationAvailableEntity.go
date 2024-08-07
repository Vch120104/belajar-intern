package transactionworkshopentities

import "time"

const TableNameWorkOrderAllocationAvailable = "trx_work_order_allocation_available"

type WorkOrderAllocationAvailable struct {
	WorkOrderAllocationAvailableId int       `gorm:"column:work_order_allocation_available_id;size:30;primaryKey" json:"work_order_allocation_available_id"`
	CompanyId                      int       `gorm:"column:company_id;size:30" json:"company_id"`
	TechnicianId                   int       `gorm:"column:technician_id;size:30" json:"technician_id"`
	ShiftCode                      string    `gorm:"column:shift_code;" json:"shift_code"`
	StartTime                      float64   `gorm:"column:start_time" json:"start_time"`
	EndTime                        float64   `gorm:"column:end_time" json:"end_time"`
	TotalHour                      float64   `gorm:"column:total_hour" json:"total_hour"`
	AvailableSystemNumber          int       `gorm:"column:available_system_number;size:30" json:"available_system_number"`
	ServiceDateTime                time.Time `gorm:"column:service_date_time" json:"service_date_time"`
	ForemanId                      int       `gorm:"column:foreman_id;size:30" json:"foreman_id"`
	ReferenceType                  string    `gorm:"column:reference_type;" json:"reference_type"`
	ReferenceSystemNumber          int       `gorm:"column:reference_system_number;size:30" json:"reference_system_number"`
	ReferenceLine                  int       `gorm:"column:reference_line;size:30" json:"reference_line"`
	Remark                         string    `gorm:"column:remark;" json:"remark"`
	CreateDate                     time.Time `gorm:"column:create_date" json:"create_date"`
	CreateBy                       int       `gorm:"column:create_by;size:30" json:"create_by"`
	ChangeDate                     time.Time `gorm:"column:change_date" json:"change_date"`
	ChangeBy                       int       `gorm:"column:change_by;size:30" json:"change_by"`
}

func (*WorkOrderAllocationAvailable) TableName() string {
	return TableNameWorkOrderAllocationAvailable
}
