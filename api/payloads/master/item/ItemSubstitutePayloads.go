package masteritempayloads

import "time"

type ItemSubstitutePayloads struct {
	ItemSubstituteId int    `json:"item_substitute_id" parent_entity:"mtr_item_substitute" main_table:"mtr_item_substitute"`
	EffectiveDate    string `json:"effective_date" parent_entity:"mtr_item_substitute"`
	SubstituteTypeId int    `json:"substitute_type_id" parent_entity:"mtr_item"`
	Description      string `json:"description" parent_entity:"mtr_item_substitute"`
	ItemGroupId      int    `json:"item_group_id" parent_entity:"mtr_item"`
	ItemId           int    `json:"item_id" parent_entity:"mtr_item" references:"mtr_item"`
	ItemCode         string `json:"item_code" parent_entity:"mtr_item"`
	ItemName         string `json:"item_name" parent_entity:"mtr_item"`
	IsActive         bool   `json:"is_active"`
	ItemClassId      int    `json:"item_class_id"`
	ItemClassCode    string `json:"item_class_code"`
}

type ItemSubstitutePostPayloads struct {
	IsActive         bool      `json:"is_active"`
	SubstituteTypeId int       `json:"substitute_type_id"`
	ItemSubstituteId int       `json:"item_substitute_id"`
	EffectiveDate    time.Time `json:"effective_date"`
	ItemId           int       `json:"item_id"`
	ItemGroupId      int       `json:"item_group_id"`
	Description      string    `json:"description"`
}

type ItemDetailForSubstitute struct {
	ItemId   int    `json:"item_id"`
	ItemName string `json:"item_name"`
}

type Itemforfilter struct {
	ItemId           int    `json:"item_id"`
	ItemCode         string `json:"item_code"`
	ItemName         string `json:"item_name"`
	ItemClassCode    string `json:"item_class_code"`
	ItemTypeCode     string `json:"item_type_code"`
	ItemLevel_1_Code string `json:"item_level_1_code"`
	ItemLevel_2_Code string `json:"item_level_2_code"`
	ItemLevel_3_Code string `json:"Item_level_3_code"`
	ItemLevel_4_Code string `json:"item_level_4_code"`
}

type ItemSubstituteCode struct {
	SubstituteTypeId    int    `json:"substitute_type_id"`
	SubstituteTypeNames string `json:"substitute_type_name"`
}
