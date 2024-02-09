package masteritempayloads

type ItemSubstituteDetailPayloads struct {
	ItemSubstituteDetailId int     `gorm:"column:item_substitute_detail_id;size:30;not null;primaryKey" json:"item_substitute_detail_id"`
	ItemSubstituteId       int     `gorm:"column:item_substitute_id;size:30;not null" json:"item_substitute_id"`
	Quantity               float64 `gorm:"column:quantity;not null" json:"quantity"`
	Sequence               int     `gorm:"column:sequence;not null;size:30" json:"sequence"`
}
