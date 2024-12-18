package masteritementities

const TableNameItemTransferStatus = "mtr_item_transfer_status"

type ItemTransferStatus struct {
	ItemTransferStatusId          int    `gorm:"column:item_transfer_status_id;size:30;not null;primaryKey" json:"item_transfer_status_id"`
	ItemTransferStatusCode        string `gorm:"column:item_transfer_status_code;size:100;not null" json:"item_transfer_status_code"`
	ItemTransferStatusDescription string `gorm:"column:item_transfer_status_description;size:50;not null" json:"item_transfer_status_description"`
	OrderNumber                   *int   `gorm:"column:order_number;size:30;null" json:"order_number"`
}

func (*ItemTransferStatus) TableName() string {
	return TableNameItemTransferStatus
}
