package masterpayloads

import "time"

type PointProspectingRequest struct {
	RecordStatus  string    `json:"record_status" validate:"required"`
	PointVariable string    `json:"point_variable" validate:"required"`
	PointValue    int       `json:"point_value" validate:"required"`
	EffectiveDate time.Time `json:"effective_date" validate:"required"`
	UserIdCreated string    `json:"user_id_created" validate:"required"`
}

type PointProspectingUpdateStatus struct {
	RecordStatus string `json:"record_status"`
}

type PointProspectingUpdateRequest struct {
	PointVariable string `json:"point_variable" validate:"required"`
	PointValue    int    `json:"point_value" validate:"required"`
	EffectiveDate string `json:"effective_date" validate:"required"`
}

type PointProspectingResponse struct {
	RecordStatus  string    `gorm:"column:RECORD_STATUS" json:"record_status"`
	PointVariable string    `gorm:"column:POINT_VARIABLE" json:"point_variable"`
	PointValue    int       `gorm:"column:POINT_VALUE" json:"point_value"`
	EffectiveDate time.Time `gorm:"column:EFFECTIVE_DATE" json:"effective_date"`
}
