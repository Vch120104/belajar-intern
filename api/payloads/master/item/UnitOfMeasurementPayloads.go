package masteritempayloads

type UomResponse struct {
	IsActive       bool   `json:"is_active" parent_entity:"mtr_uom"`
	UomId          int    `json:"uom_id" parent_entity:"mtr_uom" main_table:"mtr_uom"`
	UomTypeId      int    `json:"uom_type_id" parent_entity:"mtr_uom_type" references:"mtr_uom_type"`
	UomTypeDesc    string `json:"uom_type_desc" parent_entity:"mtr_uom_type"`
	UomCode        string `json:"uom_code" parent_entity:"mtr_uom"`
	UomDescription string `json:"uom_description" parent_entity:"mtr_uom"`
}
