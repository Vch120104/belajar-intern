package masteritempayloads

type ItemTypeRequest struct {
	ItemTypeId   int    `json:"item_type_id"`
	ItemTypeCode string `json:"item_type_code"`
	ItemTypeName string `json:"item_type_name"`
}

type ItemTypeResponse struct {
	IsActive     bool   `json:"is_active"`
	ItemTypeId   int    `json:"item_type_id"`
	ItemTypeCode string `json:"item_type_code"`
	ItemTypeName string `json:"item_type_name"`
}

type ItemTypeDropDownResponse struct {
	IsActive     bool   `json:"is_active"`
	ItemTypeId   int    `json:"item_type_id"`
	ItemTypeCode string `json:"item_type_code"`
	ItemTypeName string `json:"item_type_name"`
}
