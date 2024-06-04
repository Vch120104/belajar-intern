package masteritempayloads

type UomResponse struct {
	IsActive           bool   `json:"is_active" parent_entity:"mtr_uom"`
	UomId              int    `json:"uom_id" parent_entity:"mtr_uom" main_table:"mtr_uom"`
	UomTypeId          int    `json:"uom_type_id" parent_entity:"mtr_uom_type" references:"mtr_uom_type"`
	UomTypeDescription string `json:"uom_type_description" parent_entity:"mtr_uom_type"`
	UomCode            string `json:"uom_code" parent_entity:"mtr_uom"`
	UomDescription     string `json:"uom_description" parent_entity:"mtr_uom"`
}

type UomIdCodeResponse struct {
	IsActive       bool   `json:"is_active"`
	UomId          int    `json:"uom_id"`
	UomTypeId      int    `json:"uom_type_id"`
	UomCode        string `json:"uom_code"`
	UomDescription string `json:"uom_description"`
}
