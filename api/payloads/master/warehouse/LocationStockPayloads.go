package masterwarehousepayloads

import "time"

type LocationStockDBResponse struct {
	CompanyId              int     `parent_entity:"mtr_location_stock" json:"company_id"` //Fk with mtr_brand on sales service
	PeriodYear             string  `parent_entity:"mtr_location_stock" json:"period_year"`
	PeriodMonth            string  `parent_entity:"mtr_location_stock" json:"period_month"`
	WarehouseId            int     `parent_entity:"mtr_location_stock" json:"warehouse_id"`
	LocationId             int     `parent_entity:"mtr_location_stock" json:"location_id"`
	ItemId                 int     `parent_entity:"mtr_location_stock" json:"item_id"`
	WarehouseGroupId       int     `parent_entity:"mtr_location_stock" json:"warehouse_group_id"`
	QuantityBegin          float64 `parent_entity:"mtr_location_stock" json:"quantity_begin"`
	QuantitySales          float64 `parent_entity:"mtr_location_stock" json:"quantity_sales"`
	QuantitySalesReturn    float64 `parent_entity:"mtr_location_stock" json:"quantity_sales_return"`
	QuantityPurchase       float64 `parent_entity:"mtr_location_stock" json:"quantity_purchase"`
	QuantityPurchaseReturn float64 `parent_entity:"mtr_location_stock" json:"quantity_purchase_return"`
	QuantityTransferIn     float64 `parent_entity:"mtr_location_stock" json:"quantity_transfer_in"`
	QuantityTransferOut    float64 `parent_entity:"mtr_location_stock" json:"quantity_transfer_out"`
	QuantityClaimIn        float64 `parent_entity:"mtr_location_stock" json:"quantity_claim_in"`
	QuantityClaimOut       float64 `parent_entity:"mtr_location_stock" json:"quantity_claim_out"`
	QuantityRobbingIn      float64 `parent_entity:"mtr_location_stock" json:"quantity_robbing_in"`
	QuantityRobbingOut     float64 `parent_entity:"mtr_location_stock" json:"quantity_robbing_out"`
	QuantityAdjustment     float64 `parent_entity:"mtr_location_stock" json:"quantity_adjustment"`
	QuantityAllocated      float64 `parent_entity:"mtr_location_stock" json:"quantity_allocated"`
	QuantityInTransit      float64 `parent_entity:"mtr_location_stock" json:"quantity_in_transit"`
	QuantityEnding         float64 `parent_entity:"mtr_location_stock" json:"quantity_ending"`
	WarehouseCostingTypeId int     `parent_entity:"mtr_warehouse_master" json:"warehouse_costing_type_id"` //mtr_warehouse
	BrandId                int     `parent_entity:"mtr_warehouse_master" json:"brand_id"`
	QuantityOnHand         float64 `gorm:"column:quantity_on_hand" json:"quantity_on_hand"`
	QuantityAvailable      float64 `gorm:"column:quantity_available" json:"quantity_available"`
}
type LocationStockUpdatePayloads struct {
	CompanyId                int        `gorm:"column:company_id;not null"        json:"company_id"` //Fk with mtr_brand on sales service
	PeriodYear               string     `gorm:"column:period_year;size:4;not null"        json:"period_year"`
	PeriodMonth              string     `gorm:"column:period_month;size:2;not null"        json:"period_month"`
	WarehouseId              int        `gorm:"column:warehouse_id;not null"        json:"warehouse_id"`
	LocationId               int        `gorm:"column:location_id;not null"        json:"location_id"`
	ItemId                   int        `gorm:"column:item_id;not null"        json:"item_id"`
	WarehouseGroupId         int        `gorm:"column:warehouse_group_id;size:30;not null"        json:"warehouse_group_id"`
	QuantityBegin            float64    `gorm:"column:quantity_begin;null"        json:"quantity_begin"`
	QuantitySales            float64    `gorm:"column:quantity_sales;null"        json:"quantity_sales"`
	QuantitySalesReturn      float64    `gorm:"column:quantity_sales_return;null"        json:"quantity_sales_return"`
	QuantityPurchase         float64    `gorm:"column:quantity_purchase;null"        json:"quantity_purchase"`
	QuantityPurchaseReturn   float64    `gorm:"column:quantity_purchase_return;null"        json:"quantity_purchase_return"`
	QuantityTransferIn       float64    `gorm:"column:quantity_transfer_in;null"        json:"quantity_transfer_in"`
	QuantityTransferOut      float64    `gorm:"column:quantity_transfer_out;null"        json:"quantity_transfer_out"`
	QuantityClaimIn          float64    `gorm:"column:quantity_claim_in;null"        json:"quantity_claim_in"`
	QuantityClaimOut         float64    `gorm:"column:quantity_claim_out;null"        json:"quantity_claim_out"`
	QuantityAdjustment       float64    `gorm:"column:quantity_adjustment;null"        json:"quantity_adjustment"`
	QuantityAllocated        float64    `gorm:"column:quantity_allocated;null"        json:"quantity_allocated"`
	QuantityInTransit        float64    `gorm:"column:quantity_in_transit;null"        json:"quantity_in_transit"`
	QuantityEnding           float64    `gorm:"column:quantity_ending;null"        json:"quantity_ending"`
	QuantityRobbingIn        float64    `json:"quantity_robbing_in"`
	QuantityRobbingOut       float64    `json:"quantity_robbing_out"`
	QuantityAssemblyIn       float64    `gorm:"column:quantity_assembly_in;null"        json:"quantity_assembly_in"`
	QuantityAssemblyOut      float64    `gorm:"column:quantity_assembly_out;null"        json:"quantity_assembly_out"`
	StockTransactionTypeId   int        `json:"stock_transaction_type_id"`
	StockTransactionReasonId int        `json:"stock_transaction_reason_id"`
	CreatedByUserId          int        `gorm:"column:created_by_user_id;size:30;" json:"created_by_user_id"`
	CreatedDate              *time.Time `gorm:"column:created_date" json:"created_date"`
	UpdatedByUserId          int        `gorm:"column:updated_by_user_id;size:30;" json:"updated_by_user_id"`
	UpdatedDate              *time.Time `gorm:"column:updated_date" json:"updated_date"`
	//TransType
}
