package transactionworkshopentities

import "time"

const TableNameRecallDetail = "trx_recall_detail"

type RecallDetail struct {
	RecallDetailId       int       `gorm:"column:recall_detail_id;size:30;primaryKey" json:"recall_detail_id"`
	RecallSystemNumber   int       `gorm:"column:recall_system_number;size:30;" json:"recall_system_number"`
	RecallLineNumber     int       `gorm:"column:recall_line_number;size:30;" json:"recall_line_number"`
	BrandId              int       `gorm:"column:brand_id;" json:"brand_id"`
	ModelId              int       `gorm:"column:model_id;" json:"model_id"`
	VariantId            int       `gorm:"column:variant_id;" json:"variant_id"`
	ColourId             int       `gorm:"column:colour_id;" json:"colour_id"`
	VehicleId            int       `gorm:"column:vehicle_id;" json:"vehicle_id"`
	VehicleChassisNumber string    `gorm:"column:vechicle_chassis_number;" json:"vechicle_chassis_number"`
	VehicleEngineNumber  string    `gorm:"column:vechicle_engine_number;" json:"vechicle_engine_number"`
	HasRecall            bool      `gorm:"column:has_recall;" json:"has_recall"`
	IsActive             bool      `gorm:"column:is_active;" json:"is_active"`
	RecallBy             string    `gorm:"column:recall_by;" json:"recall_by"`
	RecallDate           time.Time `gorm:"column:recall_date;" json:"recall_date"`
}

func (*RecallDetail) TableName() string {
	return TableNameRecallDetail
}
