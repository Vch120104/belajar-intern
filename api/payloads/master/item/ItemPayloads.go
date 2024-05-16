package masteritempayloads

type ItemResponse struct {
	IsActive                     bool    `json:"is_active"`
	ItemId                       int     `json:"item_id"`
	ItemCode                     string  `json:"item_code"`
	ItemClassId                  int     `json:"item_class_id"`
	ItemName                     string  `json:"item_name"`
	ItemGroupId                  int     `json:"item_group_id"`
	ItemType                     string  `json:"item_type"`
	ItemLevel_1                  string  `json:"item_level_1"`
	ItemLevel_2                  string  `json:"item_level_2"`
	ItemLevel_3                  string  `json:"item_level_3"`
	ItemLevel_4                  string  `json:"item_level_4"`
	SupplierId                   int     `json:"supplier_id"`
	UnitOfMeasurementTypeId      int     `json:"unit_of_measurement_type_id"`
	UnitOfMeasurementSellingId   int     `json:"unit_of_measurement_selling_id"`
	UnitOfMeasurementPurchaseId  int     `json:"unit_of_measurement_purchase_id"`
	UnitOfMeasurementStockId     int     `json:"unit_of_measurement_stock_id"`
	SalesItem                    string  `json:"sales_item"`
	Lottable                     string  `json:"lottable"`
	Inspection                   string  `json:"inspection"`
	PriceListItem                string  `json:"price_list_item"`
	StockKeeping                 bool    `json:"stock_keeping"`
	DiscountId                   int     `json:"discount_id"`
	MarkupMasterId               int     `json:"markup_master_id"`
	DimensionOfLength            float64 `json:"dimension_of_length"`
	DimensionOfWidth             float64 `json:"dimension_of_width"`
	DimensionOfHeight            float64 `json:"dimension_of_height"`
	DimensionUnitOfMeasurementId int     `json:"dimension_unit_of_measurement_id"`
	Weight                       float64 `json:"weight"`
	UnitOfMeasurementWeight      string  `json:"unit_of_measurement_weight"`
	StorageTypeId                int     `json:"storage_type_id"`
	Remark                       string  `json:"remark"`
	AtpmWarrantyClaimTypeId      int     `json:"atpm_warranty_claim_type_id"`
	LastPrice                    float64 `json:"last_price"`
	UseDiscDecentralize          string  `json:"use_disc_decentralize"`
	CommonPricelist              bool    `json:"common_pricelist"`
	IsRemovable                  bool    `json:"is_removable"`
	IsMaterialPlus               bool    `json:"is_material_plus"`
	SpecialMovementId            int     `json:"special_movement_id"`
	IsItemRegulation             string  `json:"is_item_regulation"`
	IsTechnicalDefect            string  `json:"is_technical_defect"`
	IsMandatory                  bool    `json:"is_mandatory"`
	MinimumOrderQty              float64 `json:"minimum_order_qty"`
	HarmonizedNo                 string  `json:"harmonized_no"`
	AtpmSupplierId               int     `json:"atpm_supplier_id"`
	AtpmVendorSuppliability      string  `json:"atpm_vendor_suppliability"`
	PmsItem                      string  `json:"pms_item"`
	Regulation                   string  `json:"regulation"`
	AutoPickWms                  string  `json:"auto_pick_wms"`
	GmmCatalogCode               int     `json:"gmm_catalog_code"`
	PrincipalBrandParentId       int     `json:"principal_brand_parent_id"`
	ProportionalSupplyWms        string  `json:"proportional_supply_WMS"`
	Remark2                      string  `json:"remark2"`
	Remark3                      string  `json:"remark3"`
	SourceTypeId                 int     `json:"source_type_id"`
	AtpmSupplierCodeOrderId      int     `json:"atpm_supplier_code_order_id"`
	PersonInChargeId             int     `json:"person_in_charge_id"`
	SourceConvertion             float64 `json:"source_convertion"`
	TargetConvertion             float64 `json:"target_convertion"`
}

type UserDetailResponse struct {
	UserEmployeeId int    `json:"user_employee_id"`
	EmployeNo      int    `json:"employee_no"`
	EmployeeName   string `json:"employee_name"`
}

type ItemRequest struct {
	ItemId                       int     `json:"item_id"`
	ItemCode                     string  `json:"item_code"`
	ItemClassId                  int     `json:"item_class_id"`
	ItemName                     string  `json:"item_name"`
	ItemGroupId                  int     `json:"item_group_id"`
	ItemType                     string  `json:"item_type"`
	ItemLevel1                   string  `json:"item_level_1"`
	ItemLevel2                   string  `json:"item_level_2"`
	ItemLevel3                   string  `json:"item_level_3"`
	ItemLevel4                   string  `json:"item_level_4"`
	SupplierId                   int     `json:"supplier_id"`
	UnitOfMeasurementTypeId      int     `json:"unit_of_measurement_type_id"`
	UnitOfMeasurementSellingId   int     `json:"unit_of_measurement_selling_id"`
	UnitOfMeasurementPurchaseId  int     `json:"unit_of_measurement_purchase_id"`
	UnitOfMeasurementStockId     int     `json:"unit_of_measurement_stock_id"`
	SalesItem                    string  `json:"sales_item"`
	Lottable                     string  `json:"lottable"`
	Inspection                   string  `json:"inspection"`
	PriceListItem                string  `json:"price_list_item"`
	StockKeeping                 bool    `json:"stock_keeping"`
	DiscountId                   int     `json:"discount_id"`
	MarkupMasterId               int     `json:"markup_master_id"`
	DimensionOfLength            float64 `json:"dimension_of_length"`
	DimensionOfWidth             float64 `json:"dimension_of_width"`
	DimensionOfHeight            float64 `json:"dimension_of_height"`
	DimensionUnitOfMeasurementId int     `json:"dimension_unit_of_measurement_id"`
	Weight                       float64 `json:"weight"`
	UnitOfMeasurementWeight      string  `json:"unit_of_measurement_weight"`
	StorageTypeId                int     `json:"storage_type_id"`
	Remark                       string  `json:"remark"`
	AtpmWarrantyClaimTypeId      int     `json:"atpm_warranty_claim_type_id"`
	LastPrice                    float64 `json:"last_price"`
	UseDiscDecentralize          string  `json:"use_disc_decentralize"`
	CommonPricelist              bool    `json:"common_pricelist"`
	IsRemovable                  bool    `json:"is_removable"`
	IsMaterialPlus               bool    `json:"is_material_plus"`
	SpecialMovementId            int     `json:"special_movement_id"`
	IsItemRegulation             string  `json:"is_item_regulation"`
	IsTechnicalDefect            string  `json:"is_technical_defect"`
	IsMandatory                  bool    `json:"is_mandatory"`
	MinimumOrderQty              float64 `json:"minimum_order_qty"`
	HarmonizedNo                 string  `json:"harmonized_no"`
	AtpmSupplierId               int     `json:"atpm_supplier_id"`
	AtpmVendorSuppliability      string  `json:"atpm_vendor_suppliability"`
	PmsItem                      string  `json:"pms_item"`
	Regulation                   string  `json:"regulation"`
	AutoPickWms                  string  `json:"auto_pick_wms"`
	GmmCatalogCode               int     `json:"gmm_catalog_code"`
	PrincipalBrandParentId       int     `json:"principal_brand_parent_id"`
	ProportionalSupplyWms        string  `json:"proportional_supply_WMS"`
	Remark2                      string  `json:"remark2"`
	Remark3                      string  `json:"remark3"`
	SourceTypeId                 int     `json:"source_type_id"`
	AtpmSupplierCodeOrderId      int     `json:"atpm_supplier_code_order_id"`
	PersonInChargeId             int     `json:"person_in_charge_id"`
	SourceConvertion             float32 `json:"source_convertion"`
	TargetConvertion             float32 `json:"target_convertion"`
}

type AtpmOrderTypeResponse struct {
	AtpmOrderTypeDescription string `json:"atpm_order_type_description"`
	AtpmOrderTypeCode        string `json:"atpm_order_type_code"`
	AtpmOrderTypeId          int    `json:"atpm_order_type_id"`
}

type ItemLookup struct {
	IsActive    bool   `json:"is_active" parent_entity:"mtr_item"`
	ItemId      int    `json:"item_id" parent_entity:"mtr_item" main_table:"mtr_item"`
	ItemCode    string `json:"item_code" parent_entity:"mtr_item"`
	ItemName    string `json:"item_name" parent_entity:"mtr_item"`
	ItemType    string `json:"item_type" parent_entity:"mtr_item"`
	ItemGroupId int    `json:"item_group_id" parent_entity:"mtr_item"`                                                         //fk luar mtr_item_group -> item_group_name                                              // Ambil dari ItemGroupResponse
	ItemClassId int    `json:"item_class_id" parent_entity:"mtr_item_class" references:"mtr_item_class" main_table:"mtr_item"` //fk dalam item_class_id -> ItemClassName
	SupplierId  int    `json:"supplier_id" parent_entity:"mtr_item"`                                                           //fk luar mtr_supplier, supplier_code dan supplier_name
}

type UomTypeDropdownResponse struct {
	IsActive           bool   `json:"is_active"`
	UomTypeId          int    `json:"uom_type_id"`
	UomTypeDescription string `json:"uom_type_desc"`
}

type UomDropdownResponse struct {
	IsActive       bool   `json:"is_active"`
	UomId          int    `json:"uom_id"`
	UomDescription string `json:"uom_description"`
}

type ItemDetailResponse struct {
	ItemDetailId int     `json:"item_detail_id"`
	IsActive     bool    `gorm:"column:is_active" json:"is_active"`
	ItemId       int     `json:"item_id"`
	BrandId      int     `json:"brand_id"`
	ModelId      int     `json:"model_id"`
	VariantId    int     `json:"variant_id"`
	MillageEvery float64 `json:"millage_every"`
	ReturnEvery  float64 `json:"return_every"`
}

type ItemDetailRequest struct {
	ItemDetailId int     `json:"item_detail_id" parent_entity:"mtr_item_detail" main_table:"mtr_item_detail"`
	ItemId       int     `json:"item_id" parent_entity:"mtr_item_detail"`
	BrandId      int     `json:"brand_id" parent_entity:"mtr_item_detail"`
	ModelId      int     `json:"model_id" parent_entity:"mtr_item_detail"`
	VariantId    int     `json:"variant_id" parent_entity:"mtr_item_detail"`
	MillageEvery float64 `json:"millage_every" parent_entity:"mtr_item_detail"`
	ReturnEvery  float64 `json:"return_every" parent_entity:"mtr_item_detail"`
	IsActive     bool    `json:"is_active" parent_entity:"mtr_item_detail"`
}

type ItemGroupResponse struct {
	ItemGroupId   int    `json:"item_group_id"`
	ItemGroupCode string `json:"item_group_code"`
	ItemGroupName string `json:"item_group_name"`
}

type ItemClassDetailResponse struct {
	ItemClassId   int    `json:"item_class_id"`
	ItemClassCode string `json:"item_class_code"`
	ItemClassName string `json:"item_class_name"`
}

type LineTypeResponse struct {
	LineTypeId   int    `json:"line_type_id"`
	LineTypeCode string `json:"line_type_code"`
	LineTypeName string `json:"line_type_name"`
}

type SupplierMasterResponse struct {
	SupplierId   int    `json:"supplier_id"`
	SupplierCode string `json:"supplier_code"` //fk luar mtr_supplier supplier_id -> supplier_code
	SupplierName string `json:"supplier_name"` //supplier_id -> supplier_name
}

type StorageTypeResponse struct {
	StorageTypeId   int    `json:"storage_type_id"`
	StorageTypeCode string `json:"storage_type_code"`
	StorageTypeName string `json:"storage_type_name"`
}

type AtpmWarrantyClaimTypeResponse struct {
	AtpmWarrantyClaimTypeId int `json:"atpm_warranty_claim_type_id"`
}

type SpecialMovementResponse struct {
	SpecialMovementId   int    `json:"special_movement_id"`
	SpecialMovementCode string `json:"special_movement_code"`
	SpecialMovementName string `json:"special_movement_name"`
}

type WarrantyClaimTypeResponse struct {
	WarrantyClaimTypeId          int    `json:"warranty_claim_type_id"`
	WarrantyClaimTypeCode        string `json:"warranty_claim_type_code"`
	WarrantyClaimTypeDescription string `json:"warranty_claim_type_description"`
}
type AtpmSupplierResponse struct {
	AtpmSupplierId   int    `json:"supplier_id"`
	AtpmSupplierCode string `json:"supplier_code"` //fk luar mtr_supplier supplier_id -> supplier_code
	AtpmSupplierName string `json:"supplier_name"` //supplier_id -> supplier_name
}

type AtpmSupplierCodeOrderResponse struct {
	AtpmSupplierCodeOrderId   int    `json:"supplier_id"`
	AtpmSupplierCodeOrderCode string `json:"supplier_code"` //fk luar mtr_supplier supplier_id -> supplier_code
	AtpmSupplierCodeOrderName string `json:"supplier_name"` //supplier_id -> supplier_name
}

type PersonInChargeResponse struct {
	PersonInChargeId int `json:"user_id"`
}
