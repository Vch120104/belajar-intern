package masterwarehousepayloads

type SaveWarehouseMasterRequest struct {
	IsActive                      bool   `json:"is_active"`
	WarehouseCostingType          string `json:"warehouse_costing_type"`
	WarehouseKaroseri             bool   `json:"warehouse_karoseri"`
	WarehouseNegativeStock        bool   `json:"wahouse_negative_stock"`
	WarehouseReplishmentIndicator bool   `json:"warehouse_replishment_indicator"`
	WarehouseContact              string `json:"warehouse_contact"`
	WarehouseCode                 string `json:"warehouse_code"`
	AddressId                     int    `json:"address_id"`
	BrandId                       int    `json:"brand_id"`
	SupplierId                    int    `json:"supplier_id"`
	UserId                        int    `json:"user_id"`
	WarehouseSalesAllow           bool   `json:"warehouse_sales_allow"`
	WarehouseInTransit            bool   `json:"warehouse_in_transit"`
	WarehouseName                 string `json:"warehouse_name"`
	WarehouseDetailName           string `json:"warehouse_detail_name"`
	WarehouseTransitDefault       string `json:"warehouse_transit_default"`
}

type UpdateWarehouseMasterRequest struct {
	IsActive                      bool   `json:"is_active"`
	WarehouseId                   int    `json:"warehouse_id"`
	WarehouseCostingType          string `json:"warehouse_costing_type"`
	WarehouseKaroseri             bool   `json:"warehouse_karoseri"`
	WarehouseNegativeStock        bool   `json:"wahouse_negative_stock"`
	WarehouseReplishmentIndicator bool   `json:"warehouse_replishment_indicator"`
	WarehouseContact              string `json:"warehouse_contact"`
	WarehouseCode                 string `json:"warehouse_code"`
	AddressId                     int    `json:"address_id"`
	BrandId                       int    `json:"brand_id"`
	SupplierId                    int    `json:"supplier_id"`
	UserId                        int    `json:"user_id"`
	WarehouseSalesAllow           bool   `json:"warehouse_sales_allow"`
	WarehouseInTransit            bool   `json:"warehouse_in_transit"`
	WarehouseName                 string `json:"warehouse_name"`
	WarehouseDetailName           string `json:"warehouse_detail_name"`
	WarehouseTransitDefault       string `json:"warehouse_transit_default"`
}

type GetWarehouseMasterResponse struct {
	IsActive                      bool   `json:"is_active"`
	WarehouseId                   int    `json:"warehouse_id"`
	WarehouseCostingType          string `json:"warehouse_costing_type"`
	WarehouseKaroseri             bool   `json:"warehouse_karoseri"`
	WarehouseNegativeStock        bool   `json:"wahouse_negative_stock"`
	WarehouseReplishmentIndicator bool   `json:"warehouse_replishment_indicator"`
	WarehouseContact              string `json:"warehouse_contact"`
	WarehouseCode                 string `json:"warehouse_code"`
	AddressId                     int    `json:"address_id"`
	BrandId                       int    `json:"brand_id"`
	SupplierId                    int    `json:"supplier_id"`
	UserId                        int    `json:"user_id"`
	WarehouseSalesAllow           bool   `json:"warehouse_sales_allow"`
	WarehouseInTransit            bool   `json:"warehouse_in_transit"`
	WarehouseName                 string `json:"warehouse_name"`
	WarehouseDetailName           string `json:"warehouse_detail_name"`
	WarehouseTransitDefault       string `json:"warehouse_transit_default"`
}

type GetAllWarehouseMasterRequest struct {
	IsActive                      string `json:"is_active"`
	WarehouseId                   string `json:"warehouse_id"`
	WarehouseCostingType          string `json:"warehouse_costing_type"`
	WarehouseKaroseri             string `json:"warehouse_karoseri"`
	WarehouseNegativeStock        string `json:"wahouse_negative_stock"`
	WarehouseReplishmentIndicator string `json:"warehouse_replishment_indicator"`
	WarehouseContact              string `json:"warehouse_contact"`
	WarehouseCode                 string `json:"warehouse_code"`
	AddressId                     string `json:"address_id"`
	BrandId                       string `json:"brand_id"`
	SupplierId                    string `json:"supplier_id"`
	UserId                        string `json:"user_id"`
	WarehouseSalesAllow           string `json:"warehouse_sales_allow"`
	WarehouseInTransit            string `json:"warehouse_in_transit"`
	WarehouseName                 string `json:"warehouse_name"`
	WarehouseDetailName           string `json:"warehouse_detail_name"`
	WarehouseTransitDefault       string `json:"warehouse_transit_default"`
}

type GetAllWarehouseMasterResponse struct {
	IsActive                      bool                `json:"is_active"`
	WarehouseId                   int                 `json:"warehouse_id"`
	WarehouseCostingType          string              `json:"warehouse_costing_type"`
	WarehouseKaroseri             bool                `json:"warehouse_karoseri"`
	WarehouseNegativeStock        bool                `json:"wahouse_negative_stock"`
	WarehouseReplishmentIndicator bool                `json:"warehouse_replishment_indicator"`
	WarehouseContact              string              `json:"warehouse_contact"`
	WarehouseCode                 string              `json:"warehouse_code"`
	AddressId                     int                 `json:"address_id"`
	BrandId                       int                 `json:"brand_id"`
	SupplierId                    int                 `json:"supplier_id"`
	UserId                        int                 `json:"user_id"`
	WarehouseSalesAllow           bool                `json:"warehouse_sales_allow"`
	WarehouseInTransit            bool                `json:"warehouse_in_transit"`
	WarehouseName                 string              `json:"warehouse_name"`
	WarehouseDetailName           string              `json:"warehouse_detail_name"`
	WarehouseTransitDefault       string              `json:"warehouse_transit_default"`
	AddressDetails                AddressResponse     `json:"address_details"`
	BrandDetails                  BrandResponse       `json:"brand_details"`
	SupplierDetails               SupplierResponse    `json:"supplier_details"`
	UserDetails                   UserResponse        `json:"user_details"`
	JobPositionDetails            JobPositionResponse `json:"job_position_details"`
}

type GetLookupWarehouseMasterResponse struct {
	IsActive           bool   `json:"is_active"`
	WarehouseId        int    `json:"warehouse_id"`
	WarehouseName      string `json:"warehouse_name"`
	WarehouseGroupName string `json:"warehouse_group_name"`
}

type DropdownWarehouseMasterResponse struct {
	WarehouseId   int    `json:"warehouse_id"`
	WarehouseCode string `json:"warehouse_code_name"`
}

type IsActiveWarehouseMasterResponse struct {
	IsActive                      bool   `json:"is_active"`
	WarehouseId                   int    `json:"warehouse_id"`
	WarehouseCostingType          string `json:"warehouse_costing_type"`
	WarehouseKaroseri             bool   `json:"warehouse_karoseri"`
	WarehouseNegativeStock        bool   `json:"wahouse_negative_stock"`
	WarehouseReplishmentIndicator bool   `json:"warehouse_replishment_indicator"`
	WarehouseContact              string `json:"warehouse_contact"`
	WarehouseCode                 string `json:"warehouse_code"`
	AddressId                     int    `json:"address_id"`
	BrandId                       int    `json:"brand_id"`
	SupplierId                    int    `json:"supplier_id"`
	UserId                        int    `json:"user_id"`
	WarehouseSalesAllow           bool   `json:"warehouse_sales_allow"`
	WarehouseInTransit            bool   `json:"warehouse_in_transit"`
	WarehouseName                 string `json:"warehouse_name"`
	WarehouseDetailName           string `json:"warehouse_detail_name"`
	WarehouseTransitDefault       string `json:"warehouse_transit_default"`
}

type SupplierResponse struct {
	SupplierId   int    `json:"supplier_id"`
	SupplierCode string `json:"supplier_code"`
	SupplierName string `json:"supplier_name"`
}

type AddressResponse struct {
	AddressId      int    `json:"address_id"`
	AddressStreet1 string `json:"address_street_1"`
	AddressStreet2 string `json:"address_street_2"`
	AddressStreet3 string `json:"address_street_3"`
}

type BrandResponse struct {
	BrandId   int    `json:"brand_id"`
	BrandCode string `json:"brand_code"`
	BrandName string `json:"brand_name"`
}

type UserResponse struct {
	UserId        int    `json:"user_id"`
	EmployeeName  string `json:"employee_name"`
	JobPositionId int    `json:"job_position_id"`
}

type JobPositionResponse struct {
	JobPositionId   int    `json:"job_position_id"`
	JobPositionName string `json:"job_position_name"`
}
