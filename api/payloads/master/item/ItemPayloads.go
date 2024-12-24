package masteritempayloads

type ItemResponse struct {
	IsActive                     bool    `json:"is_active"`
	ItemId                       int     `json:"item_id"`
	ItemCode                     string  `json:"item_code"`
	ItemName                     string  `json:"item_name"`
	ItemClassId                  int     `json:"item_class_id"`
	ItemGroupId                  int     `json:"item_group_id"`
	ItemTypeId                   int     `json:"item_type_id"`
	Remark                       string  `json:"remark"`
	Remark2                      string  `json:"remark2"`
	Remark3                      string  `json:"remark3"`
	SourceTypeId                 *int    `json:"source_type_id"`
	SupplierId                   *int    `json:"supplier_id"`
	SupplierName                 *string `json:"supplier_name"`
	SupplierCode                 *string `json:"supplier_code"`
	AtpmSupplierId               *int    `json:"atpm_supplier_id"`
	PersonInChargeId             *int    `json:"person_in_charge_id"`
	HarmonizedNo                 string  `json:"harmonized_no"`
	MinimumOrderQty              float64 `json:"minimum_order_qty"`
	CommonPricelist              bool    `json:"common_pricelist"`
	IsRemovable                  bool    `json:"is_removable"`
	ItemLevel_1_Id               *int    `json:"item_level_1_id"`
	ItemLevel_1_Code             string  `json:"item_level_1_code"`
	ItemLevel_1_Name             string  `json:"item_level_1_name"`
	ItemLevel_2_Id               *int    `json:"item_level_2_id"`
	ItemLevel_2_Code             string  `json:"item_level_2_code"`
	ItemLevel_2_Name             string  `json:"item_level_2_name"`
	ItemLevel_3_Id               *int    `json:"item_level_3_id"`
	ItemLevel_3_Code             string  `json:"item_level_3_code"`
	ItemLevel_3_Name             string  `json:"item_level_3_name"`
	ItemLevel_4_Id               *int    `json:"item_level_4_id"`
	ItemLevel_4_Code             string  `json:"item_level_4_code"`
	ItemLevel_4_Name             string  `json:"item_level_4_name"`
	UnitOfMeasurementTypeId      *int    `json:"unit_of_measurement_type_id"`
	UnitOfMeasurementSellingId   *int    `json:"unit_of_measurement_selling_id"`
	UnitOfMeasurementPurchaseId  *int    `json:"unit_of_measurement_purchase_id"`
	UnitOfMeasurementStockId     *int    `json:"unit_of_measurement_stock_id"`
	IsSellable                   bool    `json:"is_sellable"`
	Lottable                     bool    `json:"lottable"`
	Inspection                   bool    `json:"inspection"`
	StockKeeping                 bool    `json:"stock_keeping"`
	DiscountId                   *int    `json:"discount_id"`
	MarkupMasterId               *int    `json:"markup_master_id"`
	StorageTypeId                *int    `json:"storage_type_id"`
	AtpmWarrantyClaimTypeId      *int    `json:"atpm_warranty_claim_type_id"`
	IsItemRegulation             bool    `json:"is_item_regulation"`
	ItemRegulationId             *int    `json:"item_regulation_id"`
	IsTechnicalDefect            bool    `json:"is_technical_defect"`
	PmsItem                      bool    `json:"pms_item"`
	SpecialMovementId            *int    `json:"special_movement_id"`
	AutoPickWms                  bool    `json:"auto_pick_wms"`
	ProportionalSupplyWms        bool    `json:"proportional_supply_wms"`
	PrincipalCatalogId           *int    `json:"principal_catalog_id"`
	PrincipalBrandParentId       *int    `json:"principal_brand_parent_id"`
	SourceConvertion             float64 `json:"source_convertion"`
	TargetConvertion             float64 `json:"target_convertion"`
	UomItemId                    *int    `json:"uom_item_id"`
	SourceUomId                  *int    `    json:"source_uom_id"`
	TargetUomId                  *int    `    json:"target_uom_id"`
	PriceListItem                bool    `json:"price_list_item"`
	DimensionOfLength            float64 `json:"dimension_of_length"`
	DimensionOfWidth             float64 `json:"dimension_of_width"`
	DimensionOfHeight            float64 `json:"dimension_of_height"`
	DimensionUnitOfMeasurementId *int    `json:"dimension_unit_of_measurement_id"`
	Weight                       float64 `json:"weight"`
	UnitOfMeasurementWeight      string  `json:"unit_of_measurement_weight"`
	LastPrice                    float64 `json:"last_price"`
	UseDiscDecentralize          bool    `json:"use_disc_decentralize"`
	IsMaterialPlus               bool    `json:"is_material_plus"`
	IsMandatory                  bool    `json:"is_mandatory"`
	AtpmVendorSuppliability      bool    `json:"atpm_vendor_suppliability"`
	AtpmSupplierCodeOrderId      *int    `json:"atpm_supplier_code_order_id"`
}

type UserDetailResponse struct {
	UserEmployeeId int    `json:"user_employee_id"`
	EmployeNo      int    `json:"employee_no"`
	EmployeeName   string `json:"employee_name"`
}

type LatestItemAndLineTypeResponse struct {
	ItemId     int `json:"item_id"`
	LineTypeId int `json:"line_type_id"`
}

type ItemRequest struct {
	IsActive                     bool    `json:"is_active"`
	ItemId                       int     `json:"item_id"`
	ItemCode                     string  `json:"item_code"`
	ItemClassId                  int     `json:"item_class_id"`
	ItemName                     string  `json:"item_name"`
	ItemGroupId                  int     `json:"item_group_id"`
	ItemTypeId                   int     `json:"item_type_id"`
	ItemLevel1Id                 *int    `json:"item_level_1_id"`
	ItemLevel2Id                 *int    `json:"item_level_2_id"`
	ItemLevel3Id                 *int    `json:"item_level_3_id"`
	ItemLevel4Id                 *int    `json:"item_level_4_id"`
	SupplierId                   *int    `json:"supplier_id"`
	UnitOfMeasurementTypeId      *int    `json:"unit_of_measurement_type_id"`
	UnitOfMeasurementSellingId   *int    `json:"unit_of_measurement_selling_id"`
	UnitOfMeasurementPurchaseId  *int    `json:"unit_of_measurement_purchase_id"`
	UnitOfMeasurementStockId     *int    `json:"unit_of_measurement_stock_id"`
	SourceConvertion             float64 `json:"source_convertion"`
	TargetConvertion             float64 `json:"target_convertion"`
	UomItemId                    int     `json:"uom_item_id"`
	Lottable                     bool    `json:"lottable"`
	Inspection                   bool    `json:"inspection"`
	PriceListItem                bool    `json:"price_list_item"`
	StockKeeping                 bool    `json:"stock_keeping"`
	DiscountId                   *int    `json:"discount_id"`
	MarkupMasterId               *int    `json:"markup_master_id"`
	DimensionOfLength            float64 `json:"dimension_of_length"`
	DimensionOfWidth             float64 `json:"dimension_of_width"`
	DimensionOfHeight            float64 `json:"dimension_of_height"`
	DimensionUnitOfMeasurementId *int    `json:"dimension_unit_of_measurement_id"`
	Weight                       float64 `json:"weight"`
	UnitOfMeasurementWeight      string  `json:"unit_of_measurement_weight"`
	Remark                       string  `json:"remark"`
	AtpmWarrantyClaimTypeId      *int    `json:"atpm_warranty_claim_type_id"`
	LastPrice                    float64 `json:"last_price"`
	UseDiscDecentralize          bool    `json:"use_disc_decentralize"`
	CommonPricelist              bool    `json:"common_pricelist"`
	IsRemovable                  bool    `json:"is_removable"`
	IsMaterialPlus               bool    `json:"is_material_plus"`
	SpecialMovementId            *int    `json:"special_movement_id"`
	IsItemRegulation             bool    `json:"is_item_regulation"`
	IsTechnicalDefect            bool    `json:"is_technical_defect"`
	IsMandatory                  bool    `json:"is_mandatory"`
	MinimumOrderQty              float64 `json:"minimum_order_qty"`
	HarmonizedNo                 string  `json:"harmonized_no"`
	AtpmSupplierId               *int    `json:"atpm_supplier_id"`
	AtpmVendorSuppliability      bool    `json:"atpm_vendor_suppliability"`
	PmsItem                      bool    `json:"pms_item"`
	ItemRegulationId             *int    `json:"item_regulation_id"`
	AutoPickWms                  bool    `json:"auto_pick_wms"`
	PrincipalCatalogId           *int    `json:"principal_catalog_id"`
	PrincipalBrandParentId       *int    `json:"principal_brand_parent_id"`
	ProportionalSupplyWms        bool    `json:"proportional_supply_WMS"`
	Remark2                      string  `json:"remark2"`
	Remark3                      string  `json:"remark3"`
	SourceTypeId                 *int    `json:"source_type_id"`
	AtpmSupplierCodeOrderId      *int    `json:"atpm_supplier_code_order_id"`
	PersonInChargeId             *int    `json:"person_in_charge_id"`
	IsSellable                   bool    `json:"is_sellable"`
}

type ItemSaveResponse struct {
	IsActive     bool   `json:"is_active"`
	ItemId       int    `json:"item_id"`
	ItemCode     string `json:"item_code"`
	ItemName     string `json:"item_name"`
	ItemTypeId   int    `json:"item_type_id"`
	ItemLevel1Id *int   `json:"item_level_1_id"`
	ItemLevel2Id *int   `json:"item_level_2_id"`
	ItemLevel3Id *int   `json:"item_level_3_id"`
	ItemLevel4Id *int   `json:"item_level_4_id"`
}

type AtpmOrderTypeResponse struct {
	AtpmOrderTypeDescription string `json:"atpm_order_type_description"`
	AtpmOrderTypeCode        string `json:"atpm_order_type_code"`
	AtpmOrderTypeId          int    `json:"atpm_order_type_id"`
}

type ItemLookup struct {
	IsActive      bool   `json:"is_active" parent_entity:"mtr_item"`
	ItemId        int    `json:"item_id" parent_entity:"mtr_item" main_table:"mtr_item"`
	ItemCode      string `json:"item_code" parent_entity:"mtr_item"`
	ItemName      string `json:"item_name" parent_entity:"mtr_item"`
	ItemTypeId    int    `json:"item_type_id" parent_entity:"mtr_item"`
	ItemGroupId   int    `json:"item_group_id" parent_entity:"mtr_item"`                                                         //fk luar mtr_item_group -> item_group_name                                              // Ambil dari ItemGroupResponse
	ItemClassId   int    `json:"item_class_id" parent_entity:"mtr_item_class" references:"mtr_item_class" main_table:"mtr_item"` //fk dalam item_class_id -> ItemClassName
	SupplierId    int    `json:"supplier_id" parent_entity:"mtr_item"`
	ItemClassName string `json:"item_class_name" parent_entity:"mtr_item_class" references:"mtr_item_class" main_table:"mtr_item"`
	ItemLevel_1   string `json:"item_level_1" parent_entity:"mtr_item"`
	ItemLevel_2   string `json:"item_level_2" parent_entity:"mtr_item"`
	ItemLevel_3   string `json:"item_level_3" parent_entity:"mtr_item"`
	ItemLevel_4   string `json:"item_level_4" parent_entity:"mtr_item"`
}

type UomTypeDropdownResponse struct {
	IsActive           bool   `json:"is_active"`
	UomTypeId          int    `json:"uom_type_id"`
	UomTypeDescription string `json:"uom_type_description"`
}

type UomDropdownResponse struct {
	IsActive       bool   `json:"is_active"`
	UomId          int    `json:"uom_id"`
	UomCode        string `json:"uom_code"`
	UomDescription string `json:"uom_description"`
}

type ItemDetailResponse struct {
	ItemDetailId       int     `json:"item_detail_id"`
	IsActive           bool    `gorm:"column:is_active" json:"is_active"`
	ItemId             int     `json:"item_id"`
	BrandId            int     `json:"brand_id"`
	BrandName          string  `json:"brand_name"`
	ModelId            int     `json:"model_id"`
	ModelCode          string  `json:"model_code"`
	ModelDescription   string  `json:"model_description"`
	VariantId          int     `json:"variant_id"`
	VariantCode        string  `json:"variant_code"`
	VariantDescription string  `json:"variant_description"`
	MileageEvery       float64 `json:"mileage_every"`
	ReturnEvery        float64 `json:"return_every"`
}

type GetPrincipalCatalog struct {
	PrincipalCatalogId   int    `json:"principal_catalog_id"`
	PrincipalCatalogCode string `json:"principal_catalog_code"`
}

type ItemDetailRequest struct {
	ItemDetailId int     `json:"item_detail_id" parent_entity:"mtr_item_detail" main_table:"mtr_item_detail"`
	ItemId       int     `json:"item_id" parent_entity:"mtr_item_detail"`
	BrandId      int     `json:"brand_id" parent_entity:"mtr_item_detail"`
	ModelId      int     `json:"model_id" parent_entity:"mtr_item_detail"`
	VariantId    int     `json:"variant_id" parent_entity:"mtr_item_detail"`
	MileageEvery float64 `json:"mileage_every" parent_entity:"mtr_item_detail"`
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

type SupplierMasterResponse1 struct {
	SupplierId int `json:"supplier_id"`
}

type StorageTypeResponse struct {
	StorageTypeId   int    `json:"storage_type_id"`
	StorageTypeCode string `json:"storage_type_code"`
	StorageTypeName string `json:"storage_type_name"`
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
	UserEmployeeId       int    `json:"user_employee_id"`
	PersonInChargeId     int    `json:"person_in_charge_id"`
	PersonInChargeName   string `json:"employee_name"`
	PersonInChargeNumber string `json:"employee_number"`
}

type ItemUpdateRequest struct {
	IsTechnicalDefect   bool    `json:"is_technical_defect"`
	SpecialMovementId   int     `json:"special_movement_id"`
	PrincipalCatalogId  int     `json:"principal_catalog_id"`
	UseDiscDecentralize bool    `json:"use_disc_decentralize"`
	IsSellable          bool    `json:"is_sellable"`
	SourceConvertion    float64 `json:"source_convertion"`
	TargetConvertion    float64 `json:"target_convertion"`
}

type ItemDetailUpdateRequest struct {
	MileageEvery float64 `json:"mileage_every"`
	ReturnEvery  float64 `json:"return_every"`
}

type PrincipalBrandDropdownResponse struct {
	IsActive                 bool   `json:"is_active"`
	PrincipalBrandParentId   int    `json:"principal_brand_parent_id"`
	PrincipalBrandParentCode string `json:"principal_brand_parent_code"`
}

type PrincipalBrandDropdownDescription struct {
	PrincipalBrandParentId          int    `json:"principal_brand_parent_id"`
	PrincipalBrandParentDescription string `json:"principal_brand_parent_description"`
}

type BrandModelVariantResponse struct {
	VariantId          int    `json:"variant_id"`
	VariantCode        string `json:"variant_code"`
	VariantDescription string `json:"variant_description"`
	ModelId            int    `json:"model_id"`
	ModelCode          string `json:"model_code"`
	ModelDescription   string `json:"model_description"`
	BrandId            int    `json:"brand_id"`
	BrandCode          string `json:"brand_code"`
	BrandName          string `json:"brand_name"`
}

type ItemSearch struct {
	IsActive      bool   `json:"is_active" parent_entity:"mtr_item"`
	ItemId        int    `json:"item_id" parent_entity:"mtr_item" main_table:"mtr_item"`
	ItemCode      string `json:"item_code" parent_entity:"mtr_item"`
	ItemName      string `json:"item_name" parent_entity:"mtr_item"`
	ItemTypeId    int    `json:"item_type_id" parent_entity:"mtr_item_type" references:"mtr_item_type" main_table:"mtr_item"`
	ItemTypeCode  string `json:"item_type_code" parent_entity:"mtr_item_type" references:"mtr_item_type" main_table:"mtr_item"`
	ItemGroupId   int    `json:"item_group_id" parent_entity:"mtr_item_group" references:"mtr_item_group" main_table:"mtr_item"` // fk luar mtr_item_group -> item_group_name
	ItemClassId   int    `json:"item_class_id" parent_entity:"mtr_item_class" references:"mtr_item_class" main_table:"mtr_item"` // fk dalam item_class_id -> ItemClassName
	ItemClassCode string `json:"item_class_code" parent_entity:"mtr_item_class" references:"mtr_item_class" main_table:"mtr_item"`
	SupplierId    int    `json:"supplier_id" parent_entity:"mtr_item"` // fk luar mtr_supplier, supplier_code dan supplier_name
}

type ItemInventory struct {
	ItemId        int    `json:"item_id" parent_entity:"mtr_item"`
	ItemCode      string `json:"item_code"`
	ItemName      string `json:"item_name"`
	ItemClassName string `json:"item_class_name"`
	ItemGroupCode string `json:"item_group_code"`
	ItemClassCode string `json:"item_class_code"`
	UomCode       string `json:"uom_code" parent_entity:"uom_item"`
	IsActive      bool   `json:"is_active"`
}

type ItemListTransLookUp struct {
	ItemId           int    `json:"item_id"`
	ItemCode         string `json:"item_code"`
	ItemName         string `json:"item_name"`
	ItemClassId      int    `json:"item_class_id"`
	ItemClassCode    string `json:"item_class_code"`
	ItemClassName    string `json:"item_class_name"`
	ItemTypeId       int    `json:"item_type_id"`
	ItemTypeCode     string `json:"item_type"`
	ItemLevel_1_Code string `json:"item_level_1_code"`
	ItemLevel_2_Code string `json:"item_level_2_code"`
	ItemLevel_3_Code string `json:"item_level_3_code"`
	ItemLevel_4_Code string `json:"item_level_4_code"`
}

type DeleteItemResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type CompanyResponse struct {
	CompanyId   int    `json:"company_id"`
	CompanyCode string `json:"company_code"`
	CompanyName string `json:"company_name"`
}
