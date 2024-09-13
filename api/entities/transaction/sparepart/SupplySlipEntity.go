package transactionsparepartentities

import (
	masterentities "after-sales/api/entities/master"
	"time"
)

const TableNameSupplySlip = "trx_supply_slip"

type SupplySlip struct {
	IsActive              bool       `gorm:"column:is_active;not null;default:true" json:"is_active"`
	SupplySystemNumber    int        `gorm:"column:supply_system_number;size:30;;not null;primaryKey;size:30" json:"supply_system_number"`
	SupplyDocumentNumber  string     `gorm:"column:supply_document_number;size:50;not null" json:"supply_document_number"`
	SupplyStatusId        int        `gorm:"column:supply_status_id;size:30;not null" json:"supply_status_id"`
	SupplyDate            *time.Time `gorm:"column:supply_date;not null" json:"supply_date"`
	SupplyTypeId          int        `gorm:"column:supply_type_id;size:30;;not null" json:"supply_type_id"`
	CompanyId             int        `gorm:"column:company_id;size:30;not null" json:"company_id"`
	WorkOrderSystemNumber int        `gorm:"column:work_order_system_number;size:30;not null" json:"work_order_system_number"`
	TechnicianId          int        `gorm:"column:technician_id;size:30;not null" json:"technician_id"`
	CampaignId            int        `gorm:"column:campaign_id;size:30;null" json:"campaign_id"`
	Campaign              *masterentities.CampaignMaster
	Remark                string `gorm:"column:remark;size:50;null" json:"remark"`
}

func (*SupplySlip) TableName() string {
	return TableNameSupplySlip
}
