package masterwarehouseentities

const TableNameWarehouseAuthorize = "mtr_warehouse_authorize"

type WarehouseAuthorize struct {
	WarehouseAuthorizedId int `gorm:"column:warehouse_authorize_id;size:30;not null;primaryKey;" json:"warehouse_authorize_id"`
	EmployeeId            int `gorm:"column:user_id;size:30;not null" json:"user_id"`
	CompanyId             int `gorm:"column:company_id;size:30;not null" json:"company_id"`
	WarehouseId           int `gorm:"column:warehouse_id;size:30;not null" json:"warehouse_id"`
}

func (*WarehouseAuthorize) TableName() string {
	return TableNameWarehouseAuthorize
}
