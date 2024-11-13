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

type UomItemRequest struct {
	SourceType string `json:"source_type"`
	ItemId     int    `json:"item_id"`
	Quantity   int    `json:"quantity"`
}
type UomGetQuantityConversion struct {
	SourceType string  `json:"source_type"`
	ItemId     int     `json:"item_id"`
	Quantity   float64 `json:"quantity"`
}
type UomItemResponses struct {
	SourceConvertion float64 `json:"source_convertion"`
	TargetConvertion float64 `json:"target_convertion"`
	SourceUomId      int     `json:"source_uom_id"`
	TargetUomId      int     `json:"target_uom_id"`
}
type GetQuantityConversionResponse struct {
	SourceType         string  `json:"source_type"`
	ItemId             int     `json:"item_id"`
	Quantity           float64 `json:"quantity"`
	QuantityConversion float64 `json:"quantity_conversion"`
}

//type UomGetQuantityConvertionPayloads struct {
//	SourceType int `json:"source_type_id"`
//	ItemId       int `json:"item_id"`
//	Quantity     int `json:"quantity"`
//}
