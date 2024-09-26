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
	ItemId      int    `json:"item_id"`
	ItemCode    string `json:"item_code"`
	Description string `json:"item_description"`
	ItemClass   string `json:"item_class"`
	ItemType    string `json:"item_type"`
	ItemLevel1  string `json:"item_level_1"`
	ItemLevel2  string `json:"item_level_2"`
	ItemLevel3  string `json:"Item_level_3"`
	ItemLevel4  string `json:"item_level_4"`
}

type ItemSubstituteCode struct {
	SubstituteTypeId    int    `json:"substitute_type_id"`
	SubstituteTypeNames string `json:"substitute_type_name"`
}
