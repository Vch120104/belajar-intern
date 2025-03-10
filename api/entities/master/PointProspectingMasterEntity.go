package masterentities

import "time"

type PointProspecting struct {
	RecordStatus  string    `gorm:"column:RECORD_STATUS;primary_key" json:"RECORD_STATUS"`
	PointVariable string    `gorm:"column:POINT_VARIABLE;primary_key" json:"POINT_VARIABLE"`
	PointValue    int       `gorm:"column:POINT_VALUE;primary_key" json:"POINT_VALUE"`
	EffectiveDate time.Time `gorm:"column:EFFECTIVE_DATE;primary_key" json:"EFFECTIVE_DATE"`
	NumberChanged float64   `gorm:"column:CHANGE_NO" json:"CHANGE_NO"`
	UserIdCreated string    `gorm:"column:CREATION_USER_ID" json:"CREATION_USER_ID"`
	CreatedAt     time.Time `gorm:"column:CREATION_DATETIME; autoCreateTime" json:"CREATION_DATETIME"`
	UserChanged   string    `gorm:"column:CHANGE_USER_ID" json:"CHANGE_USER_ID"`
	UpdatedAt     time.Time `gorm:"column:CHANGE_DATETIME; autoCreateTime; autoUpdateTime" json:"CHANGE_DATETIME"`
}

func (*PointProspecting) TableName() string {
	return "umPointProspecting"
}
