package masteritempayloads

type ItemGroupGetAllResponses struct {
	ItemGroupId     int    `json:"item_group_id"`
	IsActive        bool   `json:"is_active"`
	ItemGroupName   string `json:"item_group_name"`
	ItemGroupCode   string `json:"item_group_code"`
	IsItemSparepart bool   `json:"is_item_sparepart"`
}

type ItemGroupUpdatePayload struct {
	ItemGroupCode string `json:"item_group_code"`
	ItemGroupName string `json:"item_group_name"`
}

type NewItemGroupPayload struct {
	IsActive        bool   `json:"is_active"`
	ItemGroupId     int    `json:"item_group_id"`
	ItemGroupCode   string `json:"item_group_code"`
	ItemGroupName   string `json:"item_group_name"`
	IsItemSparepart bool   `json:"is_item_sparepart"`
}
