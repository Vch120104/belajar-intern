package masteritempayloads

type SaveItemLevelRequest struct {
	IsActive        bool   `json:"is_active"`
	ItemLevelId     int    `json:"item_level_id"`
	ItemLevel       string `json:"item_level"`
	ItemClassId     int    `json:"item_class_id"`
	ItemLevelParent string `json:"item_level_parent"`
	ItemLevelCode   string `json:"item_level_code"`
	ItemLevelName   string `json:"item_level_name"`
}

type GetItemLevelResponse struct {
	IsActive        bool   `json:"is_active"`
	ItemLevelId     int    `json:"item_level_id"`
	ItemLevel       string `json:"item_level"`
	ItemClassId     int    `json:"item_class_id"`
	ItemClassCode   string `json:"item_class_code"`
	ItemLevelParent string `json:"item_level_parent"`
	ItemLevelCode   string `json:"item_level_code"`
	ItemLevelName   string `json:"item_level_name"`
}

type GetItemLevelDropdownResponse struct {
	ItemLevelId     int    `json:"item_level_id"`
	ItemLevelParent string `json:"item_level_parent"`
}

type GetItemLevelResponseById struct {
	IsActive        bool   `json:"is_active"`
	ItemLevelId     int    `json:"item_level_id"`
	ItemLevel       string `json:"item_level"`
	ItemClassId     int    `json:"item_class_id"`
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
	ItemClassId     int    `json:"item_class_id"`
	ItemClassCode   string `json:"item_class_code"`
	ItemLevelParent string `json:"item_level_parent"`
}
