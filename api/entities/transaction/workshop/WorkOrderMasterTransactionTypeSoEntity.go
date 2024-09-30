package transactionworkshopentities

var CreateWorkOrderMasterTrxTypeSoTable = "mtr_work_order_transaction_type_so"

type WorkOrderMasterTrxSoType struct {
	WorkOrderTrxTypeSoId          int    `gorm:"column:transaction_type_so_id;size:30;not null;primaryKey" json:"transaction_type_so_id"`
	WorkOrderTrxTypeSoCode        string `gorm:"column:transaction_type_so_code;size:30;" json:"transaction_type_so_code"`
	WorkOrderTrxTypeSoDescription string `gorm:"column:transaction_type_so_description;size:50;" json:"transaction_type_so_description"`
	WorkOrderTrxTypeSoPrefix      string `gorm:"column:transaction_type_so_prefix;size:30;" json:"transaction_type_so_prefix"`
	IsActive                      bool   `gorm:"column:is_active;size:1;" json:"is_active"`
}

func (*WorkOrderMasterTrxSoType) TableName() string {
	return CreateWorkOrderMasterTrxTypeSoTable
}
