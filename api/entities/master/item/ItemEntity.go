package masteritementities

import masterentities "after-sales/api/entities/master"

var CreateItemTable = "mtr_item"

type Item struct {
	IsActive                     bool                                           `gorm:"column:is_active;not null" json:"is_active"`
	ItemId                       int                                            `gorm:"column:item_id;primaryKey" json:"item_id"`
	ItemCode                     string                                         `gorm:"column:item_code;size:50;unique;not null" json:"item_code"`
	ItemClassId                  int                                            `gorm:"column:item_class_id;not null" json:"item_class_id"`
	ItemName                     string                                         `gorm:"column:item_name;size:100" json:"item_name"`
	ItemGroupId                  int                                            `gorm:"column:item_group_id;not null" json:"item_group_id"`
	ItemType                     string                                         `gorm:"column:item_type;size:1" json:"item_type"`
	ItemLevel1                   string                                         `gorm:"column:item_level_1;size:10" json:"item_level_1"`
	ItemLevel2                   string                                         `gorm:"column:item_level_2;size:10" json:"item_level_2"`
	ItemLevel3                   string                                         `gorm:"column:item_level_3;size:10" json:"item_level_3"`
	ItemLevel4                   string                                         `gorm:"column:item_level_4;size:10" json:"item_level_4"`
	SupplierId                   int                                            `gorm:"column:supplier_id;not null" json:"supplier_id"`
	UnitOfMeasurementTypeId      int                                            `gorm:"column:unit_of_measurement_type_id;not null" json:"unit_of_measurement_type_id"`
	UnitOfMeasurementSellingId   int                                            `gorm:"column:unit_of_measurement_selling_id" json:"unit_of_measurement_selling_id"`
	UnitOfMeasurementPurchaseId  int                                            `gorm:"column:unit_of_measurement_purchase_id" json:"unit_of_measurement_purchase_id"`
	UnitOfMeasurementStockId     int                                            `gorm:"column:unit_of_measurement_stock_id" json:"unit_of_measurement_stock_id"`
	SalesItem                    string                                         `gorm:"column:sales_item;size:1" json:"sales_item"`
	Lottable                     bool                                           `gorm:"column:lottable" json:"lottable"`
	Inspection                   bool                                           `gorm:"column:inspection" json:"inspection"`
	PriceListItem                string                                         `gorm:"column:price_list_item;size:1" json:"price_list_item"`
	StockKeeping                 bool                                           `gorm:"column:stock_keeping;not null" json:"stock_keeping"`
	DiscountId                   int                                            `gorm:"column:discount_id;not null" json:"discount_id"`
	MarkupMasterId               int                                            `gorm:"column:markup_master_id;not null" json:"markup_master_id"`
	DimensionOfLength            float64                                        `gorm:"column:dimension_of_length" json:"dimension_of_length"`
	DimensionOfWidth             float64                                        `gorm:"column:dimension_of_width" json:"dimension_of_width"`
	DimensionOfHeight            float64                                        `gorm:"column:dimension_of_height" json:"dimension_of_height"`
	DimensionUnitOfMeasurementId int                                            `gorm:"column:dimension_unit_of_measurement_id" json:"dimension_unit_of_measurement_id"`
	Weight                       float64                                        `gorm:"column:weight" json:"weight"`
	UnitOfMeasurementWeight      string                                         `gorm:"column:unit_of_measurement_weight;size:3" json:"unit_of_measurement_weight"`
	StorageTypeId                int                                            `gorm:"column:storage_type_id;not null" json:"storage_type_id"`
	Remark                       string                                         `gorm:"column:remark;size:512" json:"remark"`
	AtpmWarrantyClaimTypeId      int                                            `gorm:"column:atpm_warranty_claim_type_id" json:"atpm_warranty_claim_type_id"`
	LastPrice                    float64                                        `gorm:"column:last_price" json:"last_price"`
	UseDiscDecentralize          string                                         `gorm:"column:use_disc_decentralize;size:1" json:"use_disc_decentralize"`
	CommonPricelist              bool                                           `gorm:"column:common_pricelist;default:false" json:"common_pricelist"`
	IsRemovable                  bool                                           `gorm:"column:is_removable" json:"is_removable"`
	IsMaterialPlus               bool                                           `gorm:"column:is_material_plus" json:"is_material_plus"`
	SpecialMovementId            int                                            `gorm:"column:special_movement_id;not null" json:"special_movement_id"`
	IsItemRegulation             bool                                           `gorm:"column:is_item_regulation" json:"is_item_regulation"`
	IsTechnicalDefect            bool                                           `gorm:"column:is_technical_defect" json:"is_technical_defect"`
	IsMandatory                  bool                                           `gorm:"column:is_mandatory" json:"is_mandatory"`
	IsSellable                   bool                                           `gorm:"column:is_sellable" json:"is_sellable"`
	IsAffiliatedTrx              bool                                           `gorm:"column:is_affiliated_trx" json:"is_affiliated_trx"`
	MinimumOrderQty              float64                                        `gorm:"column:minimum_order_qty" json:"minimum_order_qty"`
	HarmonizedNo                 string                                         `gorm:"column:harmonized_no;size:10" json:"harmonized_no"`
	AtpmSupplierId               int                                            `gorm:"column:atpm_supplier_id;not null" json:"atpm_supplier_id"`
	AtpmVendorSuppliability      bool                                           `gorm:"column:atpm_vendor_suppliability" json:"atpm_vendor_suppliability"`
	PmsItem                      bool                                           `gorm:"column:pms_item" json:"pms_item"`
	Regulation                   string                                         `gorm:"column:regulation;size:25" json:"regulation"`
	AutoPickWms                  bool                                           `gorm:"column:auto_pick_wms" json:"auto_pick_wms"`
	GmmCatalogCode               int                                            `gorm:"column:gmm_catalog_code" json:"gmm_catalog_code"`
	PrincipalBrandParentId       int                                            `gorm:"column:principal_brand_parent_id;not null" json:"principal_brand_parent_id"`
	ProportionalSupplyWms        bool                                           `gorm:"column:proportional_supply_wms" json:"proportional_supply_wms"`
	Remark2                      string                                         `gorm:"column:remark2;size:512" json:"remark2"`
	Remark3                      string                                         `gorm:"column:remark3;size:512" json:"remark3"`
	SourceTypeId                 int                                            `gorm:"column:source_type_id" json:"source_type_id"`
	AtpmSupplierCodeOrderId      int                                            `gorm:"column:atpm_supplier_code_order_id" json:"atpm_supplier_code_order_id"`
	PersonInChargeId             int                                            `gorm:"column:person_in_charge_id" json:"person_in_charge_id"`
	ItemPackageDetail            ItemPackageDetail                              `gorm:"foreignKey:ItemId;references:ItemId"`
	ItemLocation                 ItemLocation                                   `gorm:"foreignKey:ItemId;references:ItemId"`
	Bom                          Bom                                            `gorm:"foreignKey:ItemId;references:ItemId"`
	ItemClass                    *ItemClass                                     `gorm:"foreignKey:ItemClassId;references:ItemClassId"`
	ItemSubstitute               ItemSubstitute                                 `gorm:"foreignKey:ItemId;references:ItemId"`
	FieldActionItem              *masterentities.FieldActionEligibleVehicleItem `gorm:"foreignKey:ItemId;references:ItemId"`
	ItemImport                   ItemImport                                     `gorm:"foreignKey:ItemId;references:ItemId"`
	ItemDetail                   ItemDetail                                     `gorm:"foreignKey:ItemId;references:ItemId"`
}

func (*Item) TableName() string {
	return CreateItemTable
}
