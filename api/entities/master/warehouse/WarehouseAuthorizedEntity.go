package masterwarehouseentities


const TableNameWarehouseAuthorize = "mtr_warehouse_authorize"

type WarehouseAuthorize struct {
	WarehouseAuthorizedId int `gorm:"column:warehouse_authorize_id;size:30;not null;primary_key:true;autoincrement:true" json:"warehouse_authorize_id"`
	EmployeeId int `gorm:"column:employee_id;size:30;not null" json:"employee_id"`
	CompanyId int `gorm:"column:company)id;size:30;not null" jso:"company_id"`
	WarehouseId int `gorm:"column:warehouse_id;size:30;not null" json:"warehouse_id"`
}

func (*WarehouseAuthorize) TableName() string {
	return TableNameWarehouseAuthorize
}
