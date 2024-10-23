package transactionsparepartentities

import masteritementities "after-sales/api/entities/master/item"

const ItemClaimDetailTableName = "trx_item_claim_detail"

type ItemClaimDetail struct {
	ItemClaimDetailId int `gorm:"column:item_claim_detail_id;size:30;not null;primaryKey;"        json:"item_claim_detail_id"`
	ClaimSystemNumber int `gorm:"column:claim_system_number;size:30;not null"        json:"claim_system_number"`
	//ItemClaim                ItemClaim               `gorm:"foreignKey:ClaimSystemNumber"`
	ClaimLineNumber          int                     `gorm:"column:claim_line_number;not null"        json:"claim_line_number"`
	LocationItemId           int                     `gorm:"column:location_item_id;null"        json:"location_item_id"`
	ItemIds                  int                     `gorm:"column:item_id;null;size:30"        json:"item_id"`
	Item                     masteritementities.Item `gorm:"foreignKey:ItemIds;references:item_id" json:"item"`
	ItemUnitOfMeasurement    string                  `gorm:"column:item_unit_of_measurement;null;size:3"        json:"item_unit_of_measurement"`
	ItemPrice                float64                 `gorm:"column:item_price;null"        json:"item_price"`
	ItemDiscountPercentage   float64                 `gorm:"column:item_discount_percentage;null"        json:"item_discount_percentage"`
	ItemDiscountAmount       float64                 `gorm:"column:item_discount_amount;null"        json:"item_discount_amount"`
	QuantityDeliveryOrder    float64                 `gorm:"column:quantity_delivery_order;null"        json:"quantity_delivery_order"`
	QuantityShort            float64                 `gorm:"column:quantity_short;null"        json:"quantity_short"`
	QuantityDamaged          float64                 `gorm:"column:quantity_damaged;null"        json:"quantity_damaged"`
	QuantityOver             float64                 `gorm:"column:quantity_over;null"        json:"quantity_over"`
	QuantityWrong            float64                 `gorm:"column:quantity_wrong;null"        json:"quantity_wrong"`
	QuantityClaimed          float64                 `gorm:"column:quantity_claimed;null"        json:"quantity_claimed"`
	QuantityGoodsReceive     float64                 `gorm:"column:quantity_goods_receive;null"        json:"quantity_goods_receive"`
	Remark                   string                  `gorm:"column:remark;null;size:256"        json:"remark"`
	CaseNumber               string                  `gorm:"column:case_number;null;size:40"        json:"case_number"`
	GoodsReceiveSystemNumber int                     `gorm:"column:goods_receive_system_number;null"        json:"goods_receive_system_number"`
	GoodsReceiveLineNumber   int                     `gorm:"column:goods_receive_line_number;null"        json:"goods_receive_line_number"`
	QuantityBinning          float64                 `gorm:"column:quantity_binning;null"        json:"quantity_binning"`
	ItemTotal                float64                 `gorm:"column:item_total ;not null"        json:"item_total "`
	ItemTotalBaseAmount      float64                 `gorm:"column:item_total_base_amount;not null"        json:"item_total_base_amount"`
}

func (*ItemClaimDetail) TableName() string {
	return ItemClaimDetailTableName
}
