package transactionworkshopentities

import "time"

const TableNameLicenseOwnerChange = "trx_license_owner_change"

type LicenseOwnerChange struct {
	LicenseOwnerChangeId int       `gorm:"column:license_owner_change_id;primary_key;size:30" json:"license_owner_change_id"`
	BrandId              int       `gorm:"column:brand_id;not null" json:"brand_id"`
	ModelId              int       `gorm:"column:model_id;not null" json:"model_id"`
	VehicleId            int       `gorm:"column:vehicle_id;not null" json:"vehicle_id"`
	ChangeDate           time.Time `gorm:"column:change_date;not null" json:"change_date"`
	ChangeType           string    `gorm:"column:change_type;size:1" json:"change_type"`
	TnkbOld              string    `gorm:"column:vehicle_stnk_tnkb_old;size:10" json:"vehicle_stnk_tnkb_old"`
	TnkbNew              string    `gorm:"column:vehicle_stnk_tnkb_new;size:10" json:"vehicle_stnk_tnkb_new"`
	OwnerNameOld         string    `gorm:"column:vehicle_owner_name_old;size:50" json:"vehicle_owner_name_old"`
	OwnerNameNew         string    `gorm:"column:vehicle_owner_name_new;size:50" json:"vehicle_owner_name_new"`
	OwnerAddressOld      string    `gorm:"column:vehicle_owner_address_old;size:50" json:"vehicle_owner_address_old"`
	OwnerAddressNew      string    `gorm:"column:vehicle_owner_address_new;size:50" json:"vehicle_owner_address_new"`
}

func (*LicenseOwnerChange) TableName() string {
	return TableNameLicenseOwnerChange
}
