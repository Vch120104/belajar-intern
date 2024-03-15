package masteritempayloads

type GetAllItemPackageResponse struct {
	ItemPackageId   int    `json:"item_package_id" parent_entity:"mtr_item_package" main_table:"mtr_item_package"`
	IsActive        bool   `json:"is_active" parent_entity:"mtr_item_package"`
	ItemGroupId     int    `json:"item_group_id" parent_entity:"mtr_item_group" references:"mtr_item_group"`
	ItemPackageCode string `json:"item_package_code" parent_entity:"mtr_item_package"`
	ItemPackageName string `json:"item_package_name" parent_entity:"mtr_item_package"`
	ItemPackageSet  bool   `json:"item_package_set" parent_entity:"mtr_item_package"`
	Description     string `json:"description" parent_entity:"mtr_item_package"`
}

type GetItemGroupResponse struct {
	ItemGroupId   int    `json:"item_group_id"`
	ItemGroupCode string `json:"item_group_code"`
}
