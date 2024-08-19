package masterwarehousepayloads

type LocationStockDBResponse struct {
	CompanyId              int     `parent_entity:"mtr_location_stock" json:"company_id"` //Fk with mtr_brand on sales service
	PeriodYear             string  `parent_entity:"mtr_location_stock" json:"period_year"`
	PeriodMonth            string  `parent_entity:"mtr_location_stock" json:"period_month"`
	WarehouseId            int     `parent_entity:"mtr_location_stock" json:"warehouse_id"`
	LocationId             int     `parent_entity:"mtr_location_stock" json:"location_id"`
	ItemId                 int     `parent_entity:"mtr_location_stock" json:"item_id"`
	WarehouseGroup         string  `parent_entity:"mtr_location_stock" json:"warehouse_group"`
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
	WarehouseCostingType   string  `parent_entity:"mtr_warehouse_master" json:"warehouse_costing_type"` //mtr_warehouse
	BrandId                int     `parent_entity:"mtr_warehouse_master" json:"brand_id"`
	QuantityOnHand         float64 `gorm:"column:quantity_on_hand" json:"quantity_on_hand"`
	QuantityAvailable      float64 `gorm:"column:quantity_available" json:"quantity_available"`
}
