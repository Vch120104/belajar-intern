package masteritementities

var CreateMarkupRateTable = "mtr_markup_rate"

type MarkupRate struct {
	IsActive       bool         `gorm:"column:is_active;not null;default:true"        json:"is_active"`
	MarkupRateId   int          `gorm:"column:markup_rate_id;size:30;not null;primaryKey"        json:"markup_rate_id"`
	MarkupMasterId int          `gorm:"column:markup_master_id;size:30;not null"        json:"markup_master_id"`
	MarkupMaster   MarkupMaster
	OrderTypeId    int          `gorm:"column:order_type_id;size:30;not null"        json:"order_type_id"` //fk with mtr_order_type in general service
	MarkupRate     float64      `gorm:"column:markup_rate;not null"        json:"markup_rate"`
}

func (*MarkupRate) TableName() string {
	return CreateMarkupRateTable
}
