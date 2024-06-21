package masteritementities

var CreateItemTable = "mtr_item"

type Item struct {
	IsActive                     bool              `gorm:"column:is_active;type:bool;not null"        json:"is_active"`
	ItemId                       int               `gorm:"column:item_id;type:int;size:30;primaryKey"        json:"item_id"`
	ItemCode                     string            `gorm:"column:item_code;size:20;unique;not null"        json:"item_code"`
	ItemClassId                  int               `gorm:"column:item_class_id;type:int;size:30;not null"        json:"item_class_id"`
	ItemName                     string            `gorm:"column:item_name;size:100;null"        json:"item_name"`
	ItemGroupId                  int               `gorm:"column:item_group_id;type:int;size:30;not null"        json:"item_group_id"`
	ItemType                     string            `gorm:"column:item_type;size:1;null"        json:"item_type"`
	ItemLevel1                   string            `gorm:"column:item_level_1;size:10;null"        json:"item_level_1"`
	ItemLevel2                   string            `gorm:"column:item_level_2;size:10;null"        json:"item_level_2"`
	ItemLevel3                   string            `gorm:"column:item_level_3;size:10;null"        json:"item_level_3"`
	ItemLevel4                   string            `gorm:"column:item_level_4;size:10;null"        json:"item_level_4"`
	SupplierId                   int               `gorm:"column:supplier_id;type:int;size:30;not null"        json:"supplier_id"`
	UnitOfMeasurementTypeId      int               `gorm:"column:unit_of_measurement_type_id;type:int;size:30;not null"        json:"unit_of_measurement_type_id"`
	UnitOfMeasurementSellingId   int               `gorm:"column:unit_of_measurement_selling_id;type:int;size:30;null"        json:"unit_of_measurement_selling_id"`
	UnitOfMeasurementPurchaseId  int               `gorm:"column:unit_of_measurement_purchase_id;type:int;size:30;null" json:"unit_of_measurement_purchase_id"`
	UnitOfMeasurementStockId     int               `gorm:"column:unit_of_measurement_stock_id;type:int;size:30;null" json:"unit_of_measurement_stock_id"`
	SalesItem                    string            `gorm:"column:sales_item;type:char(1);null"        json:"sales_item"`
	Lottable                     string            `gorm:"column:lottable;type:char(1);null"        json:"lottable"`
	Inspection                   string            `gorm:"column:inspection;type:char(1);null"        json:"inspection"`
	PriceListItem                string            `gorm:"column:price_list_item;type:char(1);null"        json:"price_list_item"`
	StockKeeping                 bool              `gorm:"column:stock_keeping;not null"        json:"stock_keeping"`
	DiscountId                   int               `gorm:"column:discount_id;type:int;size:30;not null" json:"discount_id"`
	MarkupMasterId               int               `gorm:"column:markup_master_id;type:int;size:30;not null"        json:"markup_master_id"`
	DimensionOfLength            float64           `gorm:"column:dimension_of_length;null"        json:"dimension_of_length"`
	DimensionOfWidth             float64           `gorm:"column:dimension_of_width;null"        json:"dimension_of_width"`
	DimensionOfHeight            float64           `gorm:"column:dimension_of_height;null"        json:"dimension_of_height"`
	DimensionUnitOfMeasurementId int               `gorm:"column:dimension_unit_of_measurement_id;type:int;size:30;null" json:"dimension_unit_of_measurement_id"`
	Weight                       float64           `gorm:"column:weight;null"        json:"weight"`
	UnitOfMeasurementWeight      string            `gorm:"column:unit_of_measurement_weight;type:char(3);null"        json:"unit_of_measurement_weight"`
	StorageTypeId                int               `gorm:"column:storage_type_id;type:int;size:30;not null"        json:"storage_type_id"`
	Remark                       string            `gorm:"column:remark;size:512;null"        json:"remark"`
	AtpmWarrantyClaimTypeId      int               `gorm:"column:atpm_warranty_claim_type_id;type:int;size:30;null"        json:"atpm_warranty_claim_type_id"`
	LastPrice                    float64           `gorm:"column:last_price;null"        json:"last_price"`
	UseDiscDecentralize          string            `gorm:"column:use_disc_decentralize;type:char(1);null"        json:"use_disc_decentralize"`
	CommonPricelist              bool              `gorm:"column:common_pricelist;null"        json:"common_pricelist"`
	IsRemovable                  bool              `gorm:"column:is_removable;null"        json:"is_removable"`
	IsMaterialPlus               bool              `gorm:"column:is_material_plus;null"        json:"is_material_plus"`
	SpecialMovementId            int               `gorm:"column:special_movement_id;type:int;size:30;not null"        json:"special_movement_id"`
	IsItemRegulation             string            `gorm:"column:is_item_regulation;type:char(1);null"        json:"is_item_regulation"`
	IsTechnicalDefect            string            `gorm:"column:is_technical_defect;type:char(1);null"        json:"is_technical_defect"`
	IsMandatory                  bool              `gorm:"column:is_mandatory;null"        json:"is_mandatory"`
	IsSellable                   bool              `gorm:"column:is_sellable;null"        json:"is_sellable"`
	IsAffiliatedTrx              bool              `gorm:"column:is_affiliated_trx;null"        json:"is_affiliated_trx"`
	MinimumOrderQty              float64           `gorm:"column:minimum_order_qty;null"        json:"minimum_order_qty"`
	HarmonizedNo                 string            `gorm:"column:harmonized_no;size:10;null"        json:"harmonized_no"`
	AtpmSupplierId               int               `gorm:"column:atpm_supplier_id;type:int;size:30;not null"        json:"atpm_supplier_id"`
	AtpmVendorSuppliability      string            `gorm:"column:atpm_vendor_suppliability;type:char(1);null"        json:"atpm_vendor_suppliability"`
	PmsItem                      string            `gorm:"column:pms_item;type:char(1);null"        json:"pms_item"`
	Regulation                   string            `gorm:"column:regulation;size:25;null"        json:"regulation"`
	AutoPickWms                  string            `gorm:"column:auto_pick_wms;type:char(1);null"        json:"auto_pick_wms"`
	GmmCatalogCode               int               `gorm:"column:gmm_catalog_code;type:int;size:30;null"        json:"gmm_catalog_code"`
	PrincipalBrandParentId       int               `gorm:"column:principal_brand_parent_id;type:int;size:30;not null" json:"principal_brand_parent_id"`
	ProportionalSupplyWms        string            `gorm:"column:proportional_supply_WMS;type:char(1);null"        json:"proportional_supply_WMS"`
	Remark2                      string            `gorm:"column:remark2;size:512;null"        json:"remark2"`
	Remark3                      string            `gorm:"column:remark3;size:512;null"        json:"remark3"`
	SourceTypeId                 int               `gorm:"column:source_type_id;type:int;size:30;null"        json:"source_type_id"`
	AtpmSupplierCodeOrderId      int               `gorm:"column:atpm_supplier_code_order_id;type:int;size:30;null"        json:"atpm_supplier_code_order_id"`
	PersonInChargeId             int               `gorm:"column:person_in_charge_id;type:int;size:30;null"        json:"person_in_charge_id"`
	ItemPackageDetail            ItemPackageDetail `gorm:"foreignKey:item_id;references:item_id"`
	ItemLocation                 ItemLocation      `gorm:"foreignKey:item_id;references:item_id"`
	Bom                          Bom               `gorm:"foreignKey:item_id;references:item_id"`
	ItemClass                    *ItemClass
	ItemSubstitute               ItemSubstitute `gorm:"foreignKey:item_id;references:item_id"`
}

func (*Item) TableName() string {

	return CreateItemTable

}
