package masteritempayloads

type SaveItemLevelRequest struct {
	IsActive        bool   `json:"is_active"`
	ItemLevelId     int    `json:"item_level_id"`
	ItemLevel       int    `json:"item_level"`
	ItemClassId     int    `json:"item_class_id"`
	ItemLevelParent int    `json:"item_level_parent"`
	ItemLevelCode   string `json:"item_level_code"`
	ItemLevelName   string `json:"item_level_name"`
}

type GetItemLevelDropdownResponse struct {
	ItemLevelId       int    `json:"item_level_id"`
	ItemLevel         int    `json:"item_level"`
	ItemLevelCode     string `json:"item_level_code"`
	ItemLevelName     string `json:"item_level_name"`
	ItemLevelCodeName string `json:"item_level_code_name"`
}

type GetItemLevelResponseById struct {
	IsActive            bool   `json:"is_active"`
	ItemLevelId         int    `json:"item_level_id"`
	ItemLevel           int    `json:"item_level"`
	ItemClassId         int    `json:"item_class_id"`
	ItemLevelParentId   int    `json:"item_level_parent_id"`
	ItemLevelParentCode string `json:"item_level_parent_code"`
	ItemLevelCode       string `json:"item_level_code"`
	ItemLevelName       string `json:"item_level_name"`
}

type GetItemLevelLookUp struct {
	ItemLevel_1_Id   int    `json:"item_level_1_id"`
	ItemLevel_1_Code string `json:"item_level_1_code"`
	ItemLevel_1_Name string `json:"item_level_1_name"`
	ItemLevel_2_Id   int    `json:"item_level_2_id"`
	ItemLevel_2_Code string `json:"item_level_2_code"`
	ItemLevel_2_Name string `json:"item_level_2_name"`
	ItemLevel_3_Id   int    `json:"item_level_3_id"`
	ItemLevel_3_Code string `json:"item_level_3_code"`
	ItemLevel_3_Name string `json:"item_level_3_name"`
	ItemLevel_4_Id   int    `json:"item_level_4_id"`
	ItemLevel_4_Code string `json:"item_level_4_code"`
	ItemLevel_4_Name string `json:"item_level_4_name"`
	IsActive         string `json:"is_active"`
}

type GetAllItemLevelResponse struct {
	IsActive        bool   `json:"is_active"`
	ItemLevelId     int    `json:"item_level_id"`
	ItemLevel       string `json:"item_level"`
	ItemLevelCode   string `json:"item_level_code"`
	ItemLevelName   string `json:"item_level_name"`
	ItemClassId     int    `json:"item_class_id"`
	ItemClassCode   string `json:"item_class_code"`
	ItemLevelParent string `json:"item_level_parent"`
}
