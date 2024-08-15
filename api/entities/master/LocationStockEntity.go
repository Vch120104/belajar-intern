package masterentities

const tablename = "mtr_location_stock"

type LocationStock struct {
	ItemInquiryId          int     `gorm:"column:item_inquiry_id;not null;primaryKey" json:"item_inquiry_id"`
	CompanyId              int     `gorm:"column:company_id;not null"        json:"company_id"` //Fk with mtr_brand on sales service
	PeriodYear             string  `gorm:"column:period_year;size:4;not null"        json:"period_year"`
	PeriodMonth            string  `gorm:"column:period_month;size:2;not null"        json:"period_month"`
	WarehouseId            int     `gorm:"column:warehouse_id;not null"        json:"warehouse_id"`
	LocationId             int     `gorm:"column:location_id;not null"        json:"location_id"`
	ItemId                 int     `gorm:"column:item_id;not null"        json:"item_id"`
	WarehouseGroup         string  `gorm:"column:warehouse_group;size:50;not null"        json:"warehouse_group"`
	QuantityBegin          float64 `gorm:"column:quantity_begin;null"        json:"quantity_begin"`
	QuantitySales          float64 `gorm:"column:quantity_sales;null"        json:"quantity_sales"`
	QuantitySalesReturn    float64 `gorm:"column:quantity_sales_return;null"        json:"quantity_sales_return"`
	QuantityPurchase       float64 `gorm:"column:quantity_purchase;null"        json:"quantity_purchase"`
	QuantityTransferIn     float64 `gorm:"column:quantity_transfer_in;null"        json:"quantity_transfer_in"`
	QuantityTransferOut    float64 `gorm:"column:quantity_transfer_out;null"        json:"quantity_transfer_out"`
	QuantityClaimIn        float64 `gorm:"column:quantity_claim_in;null"        json:"quantity_claim_in"`
	QuantityClaimOut       float64 `gorm:"column:quantity_claim_out;null"        json:"quantity_claim_out"`
	QuantityAdjustment     float64 `gorm:"column:quantity_adjustment;null"        json:"quantity_adjustment"`
	QuantityAllocated      float64 `gorm:"column:quantity_allocated;null"        json:"quantity_allocated"`
	QuantityInTransit      float64 `gorm:"column:quantity_in_transit;null"        json:"quantity_in_transit"`
	QuantityEnding         float64 `gorm:"column:quantity_ending;null"        json:"quantity_ending"`
	QuantityPurchaseReturn float64 `gorm:"column:quantity_purchase_return;null"        json:"quantity_purchase_return"`
	QuantityRobbingIn      float64 `gorm:"column:quantity_robbing_in;null"        json:"quantity_robbing_in"`
	QuantityRobbingOut     float64 `gorm:"column:quantity_robbing_out;null"        json:"quantity_robbing_out"`
	QuantityAssemblyIn     float64 `gorm:"column:quantity_assembly_in;null"        json:"quantity_assembly_in"`
	QuantityAssemblyOut    float64 `gorm:"column:quantity_assembly_out;null"        json:"quantity_assembly_out"`
}

func (*LocationStock) TableName() string {
	return tablename
}
