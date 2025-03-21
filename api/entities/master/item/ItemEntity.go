package masteritementities

var CreateItemTable = "mtr_item"

type Item struct {
	IsActive                     bool                  `gorm:"column:is_active;not null" json:"is_active"`
	ItemId                       int                   `gorm:"column:item_id;size:30;primaryKey" json:"item_id"`
	ItemCode                     string                `gorm:"column:item_code;size:50;unique;not null" json:"item_code"`
	ItemClassId                  int                   `gorm:"column:item_class_id;size:30;not null" json:"item_class_id"`
	ItemName                     string                `gorm:"column:item_name;size:100" json:"item_name"`
	ItemGroupId                  int                   `gorm:"column:item_group_id;size:30;not null" json:"item_group_id"`
	ItemTypeId                   int                   `gorm:"column:item_type_id;size:30;not null" json:"item_type_id"`
	ItemLevel1Id                 *int                  `gorm:"column:item_level_1_id;size:30" json:"item_level_1_id"`
	ItemLevel2Id                 *int                  `gorm:"column:item_level_2_id;size:30" json:"item_level_2_id"`
	ItemLevel3Id                 *int                  `gorm:"column:item_level_3_id;size:30" json:"item_level_3_id"`
	ItemLevel4Id                 *int                  `gorm:"column:item_level_4_id;size:30" json:"item_level_4_id"`
	SupplierId                   *int                  `gorm:"column:supplier_id;size:30" json:"supplier_id"`
	UnitOfMeasurementTypeId      *int                  `gorm:"column:unit_of_measurement_type_id;size:30" json:"unit_of_measurement_type_id"`
	UnitOfMeasurementSellingId   *int                  `gorm:"column:unit_of_measurement_selling_id;size:30" json:"unit_of_measurement_selling_id"`
	UnitOfMeasurementPurchaseId  *int                  `gorm:"column:unit_of_measurement_purchase_id;size:30" json:"unit_of_measurement_purchase_id"`
	UnitOfMeasurementStockId     *int                  `gorm:"column:unit_of_measurement_stock_id;size:30" json:"unit_of_measurement_stock_id"`
	Lottable                     bool                  `gorm:"column:lottable" json:"lottable"`
	Inspection                   bool                  `gorm:"column:inspection" json:"inspection"`
	PriceListItem                bool                  `gorm:"column:price_list_item" json:"price_list_item"`
	StockKeeping                 bool                  `gorm:"column:stock_keeping" json:"stock_keeping"`
	DiscountId                   *int                  `gorm:"column:discount_id;size:30" json:"discount_id"`
	MarkupMasterId               *int                  `gorm:"column:markup_master_id;size:30" json:"markup_master_id"`
	DimensionOfLength            float64               `gorm:"column:dimension_of_length" json:"dimension_of_length"`
	DimensionOfWidth             float64               `gorm:"column:dimension_of_width" json:"dimension_of_width"`
	DimensionOfHeight            float64               `gorm:"column:dimension_of_height" json:"dimension_of_height"`
	DimensionUnitOfMeasurementId *int                  `gorm:"column:dimension_unit_of_measurement_id;size:30" json:"dimension_unit_of_measurement_id"`
	Weight                       float64               `gorm:"column:weight" json:"weight"`
	UnitOfMeasurementWeight      string                `gorm:"column:unit_of_measurement_weight;size:3" json:"unit_of_measurement_weight"`
	Remark                       string                `gorm:"column:remark;size:512" json:"remark"`
	AtpmWarrantyClaimTypeId      *int                  `gorm:"column:atpm_warranty_claim_type_id;size:30" json:"atpm_warranty_claim_type_id"` // fk to warranty claim type in general-service
	LastPrice                    float64               `gorm:"column:last_price" json:"last_price"`
	UseDiscDecentralize          bool                  `gorm:"column:use_disc_decentralize" json:"use_disc_decentralize"`
	CommonPricelist              bool                  `gorm:"column:common_pricelist;default:false" json:"common_pricelist"`
	IsRemovable                  bool                  `gorm:"column:is_removable" json:"is_removable"`
	IsMaterialPlus               bool                  `gorm:"column:is_material_plus" json:"is_material_plus"`
	SpecialMovementId            *int                  `gorm:"column:special_movement_id;size:30" json:"special_movement_id"` // fk to special movement in general-service
	IsItemRegulation             bool                  `gorm:"column:is_item_regulation" json:"is_item_regulation"`
	IsTechnicalDefect            bool                  `gorm:"column:is_technical_defect" json:"is_technical_defect"`
	IsMandatory                  bool                  `gorm:"column:is_mandatory" json:"is_mandatory"`
	IsSellable                   bool                  `gorm:"column:is_sellable" json:"is_sellable"`
	MinimumOrderQty              float64               `gorm:"column:minimum_order_qty" json:"minimum_order_qty"`
	HarmonizedNo                 string                `gorm:"column:harmonized_no;size:10" json:"harmonized_no"`
	AtpmSupplierId               *int                  `gorm:"column:atpm_supplier_id;size:30" json:"atpm_supplier_id"` // fk to supplier in general-service
	AtpmVendorSuppliability      bool                  `gorm:"column:atpm_vendor_suppliability" json:"atpm_vendor_suppliability"`
	PmsItem                      bool                  `gorm:"column:pms_item" json:"pms_item"`
	ItemRegulationId             *int                  `gorm:"column:item_regulation_id;size:30" json:"item_regulation_id"`
	AutoPickWms                  bool                  `gorm:"column:auto_pick_wms" json:"auto_pick_wms"`
	PrincipalCatalogId           *int                  `gorm:"column:principal_catalog_id;size:30" json:"principal_catalog_id"`
	PrincipalBrandParentId       *int                  `gorm:"column:principal_brand_parent_id;size:30" json:"principal_brand_parent_id"`
	ProportionalSupplyWms        bool                  `gorm:"column:proportional_supply_wms" json:"proportional_supply_wms"`
	Remark2                      string                `gorm:"column:remark2;size:512" json:"remark2"`
	Remark3                      string                `gorm:"column:remark3;size:512" json:"remark3"`
	SourceTypeId                 *int                  `gorm:"column:source_type_id;size:30" json:"source_type_id"`                           // fk to order type in general-service
	AtpmSupplierCodeOrderId      *int                  `gorm:"column:atpm_supplier_code_order_id;size:30" json:"atpm_supplier_code_order_id"` // fk to supplier in general-service
	PersonInChargeId             *int                  `gorm:"column:person_in_charge_id;size:30" json:"person_in_charge_id"`                 // fk to user details in general-service
	ItemClass                    ItemClass             `gorm:"foreignKey:ItemClassId;references:ItemClassId"`
	ItemType                     ItemType              `gorm:"foreignKey:ItemTypeId;references:ItemTypeId"`
	ItemGroup                    ItemGroup             `gorm:"foreignKey:ItemGroupId;references:ItemGroupId"`
	ItemLevel1                   *ItemLevel1           `gorm:"foreignKey:ItemLevel1Id;references:ItemLevel1Id"`
	ItemLevel2                   *ItemLevel2           `gorm:"foreignKey:ItemLevel2Id;references:ItemLevel2Id"`
	ItemLevel3                   *ItemLevel3           `gorm:"foreignKey:ItemLevel3Id;references:ItemLevel3Id"`
	ItemLevel4                   *ItemLevel4           `gorm:"foreignKey:ItemLevel4Id;references:ItemLevel4Id"`
	UnitOfMeasurementType        *UomType              `gorm:"foreignKey:UnitOfMeasurementTypeId;references:UomTypeId"`
	UnitOfMeasurementSelling     *Uom                  `gorm:"foreignKey:UnitOfMeasurementSellingId;references:UomId"`
	UnitOfMeasurementPurchase    *Uom                  `gorm:"foreignKey:UnitOfMeasurementPurchaseId;references:UomId"`
	UnitOfMeasurementStock       *Uom                  `gorm:"foreignKey:UnitOfMeasurementStockId;references:UomId"`
	Discount                     *Discount             `gorm:"foreignKey:DiscountId;references:DiscountCodeId"`
	MarkupMaster                 *MarkupMaster         `gorm:"foreignKey:MarkupMasterId;references:MarkupMasterId"`
	DimensionUnitOfMeasurement   *Uom                  `gorm:"foreignKey:DimensionUnitOfMeasurementId;references:UomId"`
	PrincipalCatalog             *PrincipalCatalog     `gorm:"foreignKey:PrincipalCatalogId;references:PrincipalCatalogId"`
	PrincipalBrandParent         *PrincipalBrandParent `gorm:"foreignKey:PrincipalBrandParentId;references:PrincipalBrandParentId"`
}

func (*Item) TableName() string {
	return CreateItemTable
}
