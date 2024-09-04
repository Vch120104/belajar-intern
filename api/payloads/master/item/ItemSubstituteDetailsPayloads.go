package masteritempayloads

type ItemSubstituteDetailPayloads struct {
	IsActive               bool    `json:"is_active" parent_entity:"mtr_item_substitute_detail" main_table:"mtr_item_substitute_detail"`
	ItemSubstituteDetailId int     `json:"item_substitute_detail_id" parent_entity:"mtr_item_substitute_detail" main_table:"mtr_item_substitute_detail" `
	ItemSubstituteId       int     `json:"item_substitute_id" parent_entity:"mtr_item_substitute" references:"mtr_item_substitute"`
	ItemId                 int     `json:"item_id" parent_entity:"mtr_item" references:"mtr_item"`
	ItemCode               string  `json:"item_code" parent_entity:"mtr_item"`
	ItemName               string  `json:"item_name" parent_entity:"mtr_item"`
	Quantity               float64 `json:"quantity" parent_entity:"mtr_item_substitute_detail"`
	Sequence               int     `json:"sequence" parent_entity:"mtr_item_substitute_detail"`
}

type ItemSubstituteDetailPostPayloads struct {
	ItemSubstituteDetailId int     `json:"item_substitute_detail_id"`
	ItemId                 int     `json:"item_id"`
	Quantity               float64 `json:"quantity"`
	Sequence               int     `json:"sequence"`
}

type ItemSubstituteDetailGetPayloads struct {
	IsActive               bool    `json:"is_active"`
	ItemSubstituteDetailId int     `json:"item_substitute_detail_id"`
	ItemSubstituteId       int     `json:"item_substitute_id"`
	ItemId                 int     `json:"item_id"`
	Quantity               float64 `json:"quantity"`
	Sequence               int     `json:"sequence"`
}