package transactionsparepartentities

import (
	masterwarehouseentities "after-sales/api/entities/master/warehouse"
	"time"
)

const GoodReceiveDetailTableName = "trx_goods_receive_detail"

type GoodsReceiveDetail struct {
	GoodsReceiveDetailSystemNumber int `gorm:"column:goods_receive_detail_system_number;not null;primaryKey;size:30"        json:"goods_receive_detail_system_number"`
	GoodsReceiveSystemNumber       int `gorm:"column:goods_receive_system_number;not null;size:30"        json:"goods_receive_system_number"`
	GoodsReceiveLineNumber         int `gorm:"column:goods_receive_line_number;not null;size:30"        json:"goods_receive_line_number"`
	ItemId                         int `gorm:"column:item_id;not null;size:30"        json:"item_id"`
	//Item                           masteritementities.Item
	ItemUnitOfMeasurement    string  `gorm:"column:item_unit_of_measurement;null;size:5"        json:"item_unit_of_measurement"`
	UnitOfMeasurementRate    float64 `gorm:"column:unit_of_measurement_rate;null"        json:"unit_of_measurement_rate"`
	ItemPrice                float64 `gorm:"column:item_price;null"        json:"item_price"`
	ItemDiscountPercent      float64 `gorm:"column:item_discount_percent;null"        json:"item_discount_percent"`
	ItemDiscountAmount       float64 `gorm:"column:item_discount_amount;null"        json:"item_discount_amount"`
	QuantityReference        float64 `gorm:"column:quantity_reference;null"        json:"quantity_reference"`
	QuantityDeliveryOrder    float64 `gorm:"column:quantity_delivery_order;null"        json:"quantity_delivery_order"`
	QuantityShort            float64 `gorm:"column:quantity_short;null"        json:"quantity_short"`
	QuantityDamage           float64 `gorm:"column:quantity_damage;null"        json:"quantity_damage"`
	QuantityOver             float64 `gorm:"column:quantity_over;null"        json:"quantity_over"`
	QuantityWrong            float64 `gorm:"column:quantity_wrong;null"        json:"quantity_wrong"`
	QuantityGoodsReceive     float64 `gorm:"column:quantity_goods_receive;null"        json:"quantity_goods_receive"`
	WarehouseLocationId      int     `gorm:"column:warehouse_location_id;not null;size:30"        json:"warehouse_location_id"`
	WarehouseLocation        masterwarehouseentities.WarehouseLocation
	WarehouseLocationClaimId int `gorm:"column:warehouse_location_claim_id;not null;size:30"        json:"warehouse_location_claim_id"`
	//WarehouseLocationClaim
	CaseNumber               string `gorm:"column:case_number;null"        json:"case_number;size:40"`
	BinningId                int    `gorm:"column:binning_system_number;not null;size:30"        json:"binning_system_number"`
	Binning                  BinningStock
	BinningDocumentNumber    string `gorm:"column:binning_document_number;not null;size:25"        json:"binning_document_number"`
	BinningLineNumber        int    `gorm:"column:binning_line_number;not null"        json:"binning_line_number"`
	ReferenceSystemNumber    int    `gorm:"column:reference_system_number;null;size:30"        json:"reference_system_number"`
	ReferenceLineNumber      int    `gorm:"column:reference_line_number;null"        json:"reference_line_number"`
	ClaimSystemNumber        int    `gorm:"column:claim_system_number;not null;size:30"        json:"claim_system_number"`
	ClaimLineNumber          int    `gorm:"column:claim_line_number;not null"        json:"claim_line_number"`
	InvoicePayableSystemNo   int    `gorm:"column:invoice_payable_system_no;not null"        json:"invoice_payable_system_no"`
	InvoiceReceiptLineNumber int    `gorm:"column:invoice_receipt_line_number;not null"        json:"invoice_receipt_line_number"`

	CipSystemNumber int `gorm:"column:cip_system_number;not null"        json:"cip_system_number"`
	CipLineNumber   int `gorm:"column:cip_line_number;not null"        json:"cip_line_number"`

	ItemTotal                float64 `gorm:"column:item_total;not null"        json:"item_total"`
	ItemTotalBaseAmount      float64 `gorm:"column:item_total_base_amount;not null"        json:"item_total_base_amount"`
	ItemClaimTotal           float64 `gorm:"column:item_claim_total;null"        json:"item_claim_total"`
	ItemClaimTotalBaseAmount float64 `gorm:"column:item_claim_total_base_amount;null"        json:"item_claim_total_base_amount"`
	FreightCostAmount        float64 `gorm:"column:freight_cost_amount;null"        json:"freight_cost_amount"`
	FreightCostBaseAmount    float64 `gorm:"column:freight_cost_base_amount;null"        json:"freight_cost_base_amount"`
	InsuranceCostAmount      float64 `gorm:"column:insurance_cost_amount;null"        json:"insurance_cost_amount"`
	InsuranceCostBaseAmount  float64 `gorm:"column:insurance_cost_base_amount;null"        json:"insurance_cost_base_amount"`
	OthersCostAmount         float64 `gorm:"column:others_cost_amount;null"        json:"others_cost_amount"`

	OthersCostBaseAmount      float64    `gorm:"column:others_cost_base_amount;null"        json:"others_cost_base_amount"`
	TotalLandedCost           float64    `gorm:"column:total_landed_cost;null"        json:"total_landed_cost"`
	TotalLandedCostBaseAmount float64    `gorm:"column:total_landed_cost_base_amount;null"        json:"total_landed_cost_base_amount"`
	ItemTotalPrice            float64    `gorm:"column:item_total_price;null"        json:"item_total_price"`
	ItemTotalPriceBaseAmount  float64    `gorm:"column:item_total_price_base_amount;null"        json:"item_total_price_base_amount"`
	CreatedByUserId           int        `gorm:"column:created_by_user_id;size:30;" json:"created_by_user_id"`
	CreatedDate               *time.Time `gorm:"column:created_date" json:"created_date"`
	UpdatedByUserId           int        `gorm:"column:updated_by_user_id;size:30;" json:"updated_by_user_id"`
	UpdatedDate               *time.Time `gorm:"column:updated_date" json:"updated_date"`
	ChangeNo                  int        `gorm:"column:change_no;size:30;" json:"change_no"`
}

func (*GoodsReceiveDetail) TableName() string {
	return GoodReceiveDetailTableName
}
