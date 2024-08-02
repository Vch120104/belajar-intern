package transactionworkshopentities

import "time"

const TableNameAssignTechnician = "trx_assign_technician"

type AssignTechnician struct {
	AssignTechnicianId int       `gorm:"column:assign_technician_id;size:30;primaryKey" json:"assign_technician_id"`
	CompanyId          int       `gorm:"column:company_id;size:30" json:"company_id"`
	ForemanId          int       `gorm:"column:foreman_id;size:30" json:"foreman_id"`
	ServiceDate        time.Time `gorm:"column:service_date" json:"service_date"`
	CpcCode            string    `gorm:"column:cpc_code" json:"cpc_code"`
	ShiftCode          string    `gorm:"column:shift_code" json:"shift_code"`
	TechnicianId       int       `gorm:"column:technician_id;size:30" json:"technician_id"`
	CreateDate         time.Time `gorm:"column:create_date" json:"create_date"`
	CreateBy           int       `gorm:"column:create_by;size:30" json:"create_by"`
	ChangeDate         time.Time `gorm:"column:change_date" json:"change_date"`
	ChangeBy           int       `gorm:"column:change_by;size:30" json:"change_by"`
}

func (*AssignTechnician) TableName() string {
	return TableNameAssignTechnician
}
