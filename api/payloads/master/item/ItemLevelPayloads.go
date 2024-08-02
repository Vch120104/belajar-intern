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
	ItemLevelId   int    `json:"item_level_id"`
	ItemLevelName string `json:"item_level_name"`
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

type GetItemLevelLookUp struct {
	Item_level_1      string `json:"item_level_1"`
	Item_level_1_name string `json:"item_level_1_name"`
	Item_level_2      string `json:"item_level_2"`
	Item_level_2_name string `json:"item_level_2_name"`
	Item_level_3      string `json:"item_level_3"`
	Item_level_3_name string `json:"item_level_3_name"`
	Item_level_4      string `json:"item_level_4"`
	Item_level_4_name string `json:"item_level_4_name"`
	ItemLevelId       int    `json:"item_level_id"`
	IsActive          string `json:"is_active"`
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
