package masteritempayloads

import "time"

type ItemSubstitutePayloads struct {
	BrandId          int    `json:"brand_id" parent_entity:"mtr_item_substitute"`
	SubstituteTypeId int    `json:"substitute_type_Id" parent_entity:"mtr_item_substitute"`
	ItemSubstituteId int    `json:"item_substitute_id" parent_entity:"mtr_item_substitute" main_table:"mtr_item_substitute"`
	EffectiveDate    string `json:"effective_date" parent_entity:"mtr_item_substitute"`
	ItemId           int    `json:"item_id" parent_entity:"mtr_item" references:"mtr_item"`
	ItemCode         string `json:"item_code" parent_entity:"mtr_item"`
	ItemName         string `json:"item_name" parent_entity:"mtr_item"`
	ItemGroupId      int    `json:"item_group_id" parent_entity:"mtr_item_substitute"`
	ItemClassId      int    `json:"item_class_id" parent_entity:"mtr_item_class"`
	ItemClassCode    string `json:"mtr_class_code" parent_entity:"mtr_item_class"`
}

type ItemSubstitutePayloadsGetAll struct {
	SubstituteTypeId int    `json:"substitute_type_Id" parent_entity:"mtr_item_substitute"`
	ItemSubstituteId int    `json:"item_substitute_id" parent_entity:"mtr_item_substitute" main_table:"mtr_item_substitute"`
	EffectiveDate    string `json:"effective_date" parent_entity:"mtr_item_substitute"`
	ItemId           int    `json:"item_id" parent_entity:"mtr_item" references:"mtr_item"`
	ItemCode         string `json:"item_code" parent_entity:"mtr_item"`
	ItemName         string `json:"item_name" parent_entity:"mtr_item"`
}

type ItemSubstituteById struct{
	SubstituteTypeId int    `json:"substitute_type_Id" parent_entity:"mtr_item_substitute"`
	ItemSubstituteId int    `json:"item_substitute_id" parent_entity:"mtr_item_substitute" main_table:"mtr_item_substitute"`
	EffectiveDate    string `json:"effective_date" parent_entity:"mtr_item_substitute"`
	ItemId           int    `json:"item_id" parent_entity:"mtr_item" references:"mtr_item"`
	ItemCode         string `json:"item_code" parent_entity:"mtr_item"`
	ItemName         string `json:"item_name" parent_entity:"mtr_item"`
	ItemGroupId      int    `json:"item_group_id" parent_entity:"mtr_item_substitute"`
	ItemClassId      int    `json:"item_class_id" parent_entity:"mtr_item_class"`
	ItemClassCode    string `json:"mtr_class_code" parent_entity:"mtr_item_class"`
}

type GetSubstitutePayloads struct {
	SubstituteTypeCode string `json:"substitute_type_code"`
	SubstituteTypeId   int    `json:"substitute_type_Id"`
	SubstituteTypeName string `json:"substitute_type_name"`
	IsActive           bool   `json:"is_active"`
}

type ItemSubstitutePostPayloads struct {
	SubstituteTypeId int       `json:"substitute_type_id"`
	ItemSubstituteId int       `json:"item_substitute_id"`
	EffectiveDate    time.Time `json:"effective_date"`
	ItemId           int       `json:"item_id"`
	ItemGroupId      int       `json:"item_group_id"`
	ItemClassId      int       `json:"item_class_id"`
	Description      string    `json:"description"`
}

type ItemDetailForSubstitute struct {
	ItemId   int    `json:"item_id"`
	ItemName string `json:"item_name"`
}