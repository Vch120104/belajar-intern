package transactionsparepartentities

const TableName = "mtr_purchase_order_limit"

type PurchaseOrderLimit struct {
	PurchaseOrderLimitId int `gorm:"column:purchase_order_limit_id;size:30;not null;primaryKey;" json:"purchase_order_limit_id"`
	CompanyId            int `gorm:"column:company_id;size:30;" json:"company_id"`
	OrderTypeId          int `gorm:"column:order_type_id;size:30;" json:"order_type_id"`
	OrderLimit           int `gorm:"column:order_limit;size:30;" json:"order_limit"`
}

func (*PurchaseOrderLimit) TableName() string {
	return TableName
}
