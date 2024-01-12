package masteritementities

import "after-sales/api/entities/master"

var CreateItemTable = "mtr_item"

type Item struct {
	IsActive                     bool                 `gorm:"column:is_active;size:1;not"        json:"is_active"`
	ItemId                       int                  `gorm:"column:item_id;not null;primaryKey"        json:"item_id"`
	ItemCode                     string               `gorm:"column:item_code;unique;size:20;not null"        json:"item_code"`
	ItemClassId                  int                  `gorm:"column:item_class_id;not null"        json:"item_class_id"` //FK dalam mtr_item_class
	ItemClass                    ItemClass            
	ItemName                     string               `gorm:"column:item_name;size:100;null"        json:"item_name"`
	ItemGroupId                  int                  `gorm:"column:item_group_id;not null"        json:"item_group_id"` //FK Luar with mtr_item_group common-general service
	ItemType                     string               `gorm:"column:item_type;size:1;type:char(1);null"        json:"item_type"`
	ItemLevel1                   string               `gorm:"column:item_level_1;size:10;null"        json:"item_level_1"`
	ItemLevel2                   string               `gorm:"column:item_level_2;size:10;null"        json:"item_level_2"`
	ItemLevel3                   string               `gorm:"column:item_level_3;size:10;null"        json:"item_level_3"`
	ItemLevel4                   string               `gorm:"column:item_level_4;size:10;null"        json:"item_level_4"`
	SupplierId                   int                  `gorm:"column:supplier_id;not null"        json:"supplier_id"`                                 //FK luar with mtr_supplier general service
	UnitOfMeasurementTypeId      int                  `gorm:"column:unit_of_measurement_type_id;not null"        json:"unit_of_measurement_type_id"` //FK luar with mtr_unit_of_measurement_type
	UnitOfMeasurementSellingId   int                  `gorm:"column:unit_of_measurement_selling_id;null"        json:"unit_of_measurement_selling_id"`
	UnitOfMeasurementSelling     Uom                  `gorm:"foreignKey:UnitOfMeasurementSellingId"`
	UnitOfMeasurementPurchaseId  int                  `gorm:"column:unit_of_measurement_purchase_id;null" json:"unit_of_measurement_purchase_id"`
	UnitOfMeasurementPurchase    Uom                  `gorm:"foreignKey:UnitOfMeasurementPurchaseId"`
	UnitOfMeasurementStockId     int                  `gorm:"column:unit_of_measurement_stock_id;null" json:"unit_of_measurement_stock_id"`
	UnitOfMeasurementStock       Uom                  `gorm:"foreignKey:UnitOfMeasurementStockId"`
	SalesItem                    string               `gorm:"column:sales_item;size:1;type:char(1);null"        json:"sales_item"`
	Lottable                     string               `gorm:"column:lottable;size:1;type:char(1);null"        json:"lottable"`
	Inspection                   string               `gorm:"column:inspection;size:1;type:char(1);null"        json:"inspection"`
	PriceListItem                string               `gorm:"column:price_list_item;size:1;type:char(1);null"        json:"price_list_item"`
	StockKeeping                 bool                 `gorm:"column:stock_keeping;size:1;null"        json:"stock_keeping"`
	DiscountId                   int                  `gorm:"column:discount_id;not null" json:"discount_id"`
	Discount                     masterentities.Discount             
	MarkupMasterId               int                  `gorm:"column:markup_master_id;not null"        json:"markup_master_id"` //FK dalam mtr_markup_master
	MarkupMaster                 MarkupMaster         
	DimensionOfLength            float64              `gorm:"column:dimension_of_length;null"        json:"dimension_of_length"`
	DimensionOfWidth             float64              `gorm:"column:dimension_of_width;null"        json:"dimension_of_width"`
	DimensionOfHeight            float64              `gorm:"column:dimension_of_height;null"        json:"dimension_of_height"`
	DimensionUnitOfMeasurementId int                  `gorm:"column:dimension_unit_of_measurement_id;null" json:"dimension_unit_of_measurement_id"`
	DimensionUnitOfMeasurement   Uom                  `gorm:"foreignKey:DimensionUnitOfMeasurementId"`
	Weight                       float64              `gorm:"column:weight;null"        json:"weight"`
	UnitOfMeasurementWeight      string               `gorm:"column:unit_of_measurement_weight;size:3;type:char(3);null"        json:"unit_of_measurement_weight"`
	StorageTypeId                int                  `gorm:"column:storage_type_id;not null"        json:"storage_type_id"` //FK luar with storage_type general service
	Remark                       string               `gorm:"column:remark;size:512;null"        json:"remark"`
	AtpmWarrantyClaimTypeId      int                  `gorm:"column:atpm_warranty_claim_type_id;null"        json:"atpm_warranty_claim_type_id"` //FK luar with mtr_warranty_claim_type common service
	LastPrice                    float64              `gorm:"column:last_price;null"        json:"last_price"`
	UseDiscDecentralize          string               `gorm:"column:use_disc_decentralize;size:1;type:char(1);null"        json:"use_disc_decentralize"`
	CommonPricelist              bool                 `gorm:"column:common_pricelist;size:1;null"        json:"common_pricelist"`
	IsRemovable                  bool                 `gorm:"column:is_removable;size:1;null"        json:"is_removable"`
	IsMaterialPlus               bool                 `gorm:"column:is_material_plus;size:1;null"        json:"is_material_plus"`
	SpecialMovementId            int                  `gorm:"column:special_movement_id;not null"        json:"special_movement_id"` //FK luar with mtr_special_movement common service
	IsItemRegulation             string               `gorm:"column:is_item_regulation;size:1;type:char(1);null"        json:"is_item_regulation"`
	IsTechnicalDefect            string               `gorm:"column:is_technical_defect;size:1;type:char(1);null"        json:"is_technical_defect"`
	IsMandatory                  bool                 `gorm:"column:is_mandatory;size:1;null"        json:"is_mandatory"`
	MinimumOrderQty              float64              `gorm:"column:minimum_order_qty;null"        json:"minimum_order_qty"`
	HarmonizedNo                 string               `gorm:"column:harmonized_no;size:10;null"        json:"harmonized_no"`
	AtpmSupplierId               int                  `gorm:"column:atpm_supplier_id;not null"        json:"atpm_supplier_id"` //FK luar with mtr_supplier general service
	AtpmVendorSuppliability      string               `gorm:"column:atpm_vendor_suppliability;size:1;type:char(1);null"        json:"atpm_vendor_suppliability"`
	PmsItem                      string               `gorm:"column:pms_item;size:1;type:char(1);null"        json:"pms_item"`
	Regulation                   string               `gorm:"column:regulation;size:25;null"        json:"regulation"`
	AutoPickWms                  string               `gorm:"column:auto_pick_wms;size:1;type:char(1);null"        json:"auto_pick_wms"`
	GmmCatalogCode               int                  `gorm:"column:gmm_catalog_code;null"        json:"gmm_catalog_code"`
	PrincipalBrandParentId       int                  `gorm:"column:principal_brand_parent_id;not null" json:"principal_brand_parent_id"`
	PrincipleBrandParents        PrincipleBrandParent `gorm:"foreignKey:PrincipalBrandParentId"`
	ProportionalSupplyWms        string               `gorm:"column:proportional_supply_WMS;size:1;type:char(1);null"        json:"proportional_supply_WMS"`
	Remark2                      string               `gorm:"column:remark2;size:512;null"        json:"remark2"`
	Remark3                      string               `gorm:"column:remark3;size:512;null"        json:"remark3"`
	SourceTypeId                 int                  `gorm:"column:source_type_id;null"        json:"source_type_id"`                           //fk luar with mtr_atpm_order_type common service
	AtpmSupplierCodeOrderId      int                  `gorm:"column:atpm_supplier_code_order_id;null"        json:"atpm_supplier_code_order_id"` //FK luar with mtr_supplier general service
	PersonInChargeId             int                  `gorm:"column:person_in_charge_id;null"        json:"person_in_charge_id"`                 //FK luar with mtr_user_details general service
}

func (*Item) TableName() string {

	return CreateItemTable

}
