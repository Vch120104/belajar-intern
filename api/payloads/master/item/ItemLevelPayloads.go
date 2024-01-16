package masteritempayloads

type SaveItemLevelRequest struct {
	IsActive        bool   `json:"is_active"`
	ItemLevelId     int    `json:"item_level_id"`
	ItemLevel       string `json:"item_level"`
	ItemClassCode   string `json:"item_class_code"`
	ItemLevelParent string `json:"item_level_parent"`
	ItemLevelCode   string `json:"item_level_code"`
	ItemLevelName   string `json:"item_level_name"`
}

type GetItemLevelResponse struct {
	IsActive        bool   `json:"is_active"`
	ItemLevelId     int    `json:"item_level_id"`
	ItemLevel       string `json:"item_level"`
	ItemClassCode   string `json:"item_class_code"`
	ItemLevelParent string `json:"item_level_parent"`
	ItemLevelCode   string `json:"item_level_code"`
	ItemLevelName   string `json:"item_level_name"`
}

type GetAllItemLevelResponse struct {
	IsActive        string `json:"is_active"`
	ItemLevelId     int    `json:"item_level_id"`
	ItemLevel       string `json:"item_level"`
	ItemLevelCode   string `json:"item_level_code"`
	ItemLevelName   string `json:"item_level_name"`
	ItemClassCode   string `json:"item_class_code"`
	ItemLevelParent string `json:"item_level_parent"`
}
