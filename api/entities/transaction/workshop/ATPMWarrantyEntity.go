package transactionworkshopentities

import "time"

const TableNameAtpmWarranty = "trx_atpm_warranty"

type AtpmWarranty struct {
	WarrantySystemNumber int       `gorm:"column:warranty_system_number;size:30;primaryKey" json:"warranty_system_number"`
	IsActive             bool      `gorm:"column:is_active" json:"is_active"`
	BrandId              int       `gorm:"column:brand_id;size:30" json:"brand_id"`
	ModelId              int       `gorm:"column:model_id;size:30" json:"model_id"`
	VariantId            int       `gorm:"column:variant_id;size:30" json:"variant_id"`
	MasterTypeId         int       `gorm:"column:master_type_id;size:30" json:"master_type_id"`
	FspCategoryId        int       `gorm:"column:fsp_category_id;size:30" json:"fsp_category_id"`
	EffectiveDate        time.Time `gorm:"column:effective_date" json:"effective_date"`
	ExpireMileage        int       `gorm:"column:expire_mileage;size:30" json:"expire_mileage"`
	ExpireMonth          int       `gorm:"column:expire_month;size:30" json:"expire_month"`
	TotalAfterDiscount   float64   `gorm:"column:total_after_discount" json:"total_after_discount"`
	AfsArea              string    `gorm:"column:afs_area" json:"afs_area"`
	LabourSellingPrice   float64   `gorm:"column:labour_selling_price" json:"labour_selling_price"`
	TotalFrtQty          float64   `gorm:"column:total_frt_qty" json:"total_frt_qty"`
	TotalLabour          float64   `gorm:"column:total_labour" json:"total_labour"`
	TotalPart            float64   `gorm:"column:total_part" json:"total_part"`
	RetailDateFrom       time.Time `gorm:"column:retail_date_from" json:"retail_date_from"`
	RetailDateTo         time.Time `gorm:"column:retail_date_to" json:"retail_date_to"`
}

func (*AtpmWarranty) TableName() string {
	return TableNameAtpmWarranty
}
