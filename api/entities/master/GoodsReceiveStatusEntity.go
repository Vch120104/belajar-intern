package masterentities

const GoodsReceiveStatusMasterTable = "mtr_item_goods_receive_status"

type GoodsReceiveDocumentStatus struct {
	ItemGoodsReceiveStatusId          int    `gorm:"column:item_goods_receive_status_id;not null;primaryKey;size:30"        json:"item_goods_receive_status_id"`
	ItemGoodsReceiveStatusCode        string `gorm:"column:item_goods_receive_status_code;not null"        json:"item_goods_receive_status_code"`
	ItemGoodsReceiveStatusDescription string `gorm:"column:item_goods_receive_status_description;not null"        json:"item_goods_receive_status_description"`
}

func (*GoodsReceiveDocumentStatus) TableName() string {
	return GoodsReceiveStatusMasterTable
}
