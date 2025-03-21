package masterwarehouseentities

const TableNameWarehouseMaster = "mtr_warehouse_master"

type WarehouseMaster struct {
	CompanyId                     int   `gorm:"column:company_id;size:30;not null" json:"company_id"`
	IsActive                      *bool `gorm:"column:is_active;default:true;not null" json:"is_active"`
	WarehouseCostingTypeId        int   `gorm:"column:warehouse_costing_type_id;null;type:varchar(50)" json:"warehouse_costing_type_id"`
	WarehouseCostingType          WarehouseCostingType
	WarehouseKaroseri             *bool              `gorm:"column:warehouse_karoseri;default:false;not null" json:"warehouse_karoseri"`
	WarehouseNegativeStock        *bool              `gorm:"column:warehouse_negative_stock;default:false;not null" json:"warehouse_negative_stock"`
	WarehouseReplishmentIndicator *bool              `gorm:"column:warehouse_replishment_indicator;default:false;not null" json:"warehouse_replishment_indicator"`
	WarehouseGroupId              int                `gorm:"column:warehouse_group_id;size:30;not null" json:"warehouse_group_id"`
	WarehouseContact              string             `gorm:"column:warehouse_contact;not null;type:varchar(100)" json:"warehouse_contact"`
	WarehouseCode                 string             `gorm:"column:warehouse_code;not null;type:varchar(5);unique" json:"warehouse_code"`
	WarehouseId                   int                `gorm:"column:warehouse_id;size:30;not null;primaryKey" json:"warehouse_id"`
	AddressId                     int                `gorm:"column:address_id;size:30;not null" json:"address_id"`
	BrandId                       int                `gorm:"column:brand_id;size:30;not null" json:"brand_id"`
	SupplierId                    int                `gorm:"column:supplier_id;size:30;not null" json:"supplier_id"`
	UserId                        int                `gorm:"column:user_id;size:30;not null" json:"user_id"`
	WarehouseSalesAllow           *bool              `gorm:"column:warehouse_sales_allow;default:false;not null" json:"warehouse_sales_allow"`
	WarehouseInTransit            *bool              `gorm:"column:warehouse_in_transit;default:false;not null" json:"warehouse_in_transit"`
	WarehouseName                 string             `gorm:"column:warehouse_name;not null;type:varchar(100)" json:"warehouse_name"`
	WarehouseDetailName           string             `gorm:"column:warehouse_detail_name;not null;type:varchar(100)" json:"warehouse_detail_name"`
	WarehouseTransitDefault       string             `gorm:"column:warehouse_transit_default;not null;type:varchar(5)" json:"warehouse_transit_default"`
	WarehousePhoneNumber          string             `gorm:"column:warehouse_phone_number;not null;size:30;default:'-'" json:"warehouse_phone_number"`
	WarehouseFaxNumber            string             `gorm:"column:warehouse_fax_number;size:30" json:"warehouse_fax_number"`
	WarehouseAuthorized           WarehouseAuthorize `gorm:"foreignkey:warehouse_id;references:warehouse_id" json:"warehouse_authorized"`
}

func (*WarehouseMaster) TableName() string {
	return TableNameWarehouseMaster
}
