package transactionworkshopentities

const TableNameBookingAllocation = "trx_booking_allocation"

type BookingAllocation struct {
	BookingAllocationSystemNumber int     `gorm:"column:booking_allocation_system_number;size:30;primaryKey" json:"booking_allocation_system_number"`
	IsActive                      bool    `gorm:"column:is_active;size:1" json:"is_active"`
	CompanyId                     int     `gorm:"column:company_id;size:30" json:"company_id"`
	BrandId                       int     `gorm:"column:brand_id;size:30" json:"brand_id"`
	ShiftCode                     string  `gorm:"column:shift_code;size:30" json:"shift_code"`
	BookingSystemNumber           int     `gorm:"column:booking_system_number;size:30" json:"booking_system_number"`
	VehicleId                     int     `gorm:"column:vehicle_id;size:30" json:"vehicle_id"`
	TechnicianId                  int     `gorm:"column:technician_id;size:30" json:"technician_id"`
	BookingAllocationDate         string  `gorm:"column:booking_allocation_date" json:"booking_allocation_date"`
	BookingAllocationStartTime    float64 `gorm:"column:booking_allocation_start_time" json:"booking_allocation_start_time"`
	BookingAllocationEndTime      float64 `gorm:"column:booking_allocation_end_time" json:"booking_allocation_end_time"`
	BookingAllocationTotalHour    float64 `gorm:"column:booking_allocation_total_hour" json:"booking_allocation_total_hour"`
	BookingAllocationTechnician   int     `gorm:"column:booking_allocation_technician" json:"booking_allocation_technician"`
}

func (*BookingAllocation) TableName() string {
	return TableNameBookingAllocation
}
