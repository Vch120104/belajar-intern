package masterentities

var GoodsReceiveMasterTable = "mtr_reference_type_goods_receive"

type GoodsReceiveReferenceType struct {
	ReferenceTypeGoodReceiveId           int    `gorm:"column:reference_type_good_receive_id;not null;primaryKey;size:30"        json:"reference_type_good_receive_id"`
	ReferenceTypeGoodReceiveCode         string `gorm:"column:reference_type_good_receive_code;not null;size:25"        json:"reference_type_good_receive_code"`
	ReferenceTypeGoodsReceiveDescription string `gorm:"column:reference_type_goods_receive_description;not null;size:50"        json:"reference_type_goods_receive_description"`
}

func (*GoodsReceiveReferenceType) TableName() string {
	return GoodsReceiveMasterTable
}
