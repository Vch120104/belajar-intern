package masteritempayloads

type ItemPackageDetailResponse struct {
	ItemPackageDetailId int     `json:"item_package_detail_id" parent_entity:"mtr_item_package_detail" main_table:"mtr_item_package_detail"`
	IsActive            bool    `json:"is_active" parent_entity:"mtr_item_package_detail"`
	ItemPackageId       int     `json:"item_package_id" parent_entity:"mtr_item_package_detail" references:"mtr_item_package"`
	ItemId              int     `json:"item_id" parent_entity:"mtr_item"  references:"mtr_item"`
	ItemCode            string  `json:"item_code" parent_entity:"mtr_item" `
	ItemName            string  `json:"item_name" parent_entity:"mtr_item" `
	ItemClassId         int     `json:"item_class_id" parent_entity:"mtr_item_package"  references:"mtr_item_class"`
	ItemClassCode       string  `json:"item_class_code" parent_entity:"mtr_item_class" `
	Quantity            float64 `json:"quantity" parent_entity:"mtr_item_package_detail"`
}

type SaveItemPackageDetail struct {
	IsActive            bool    `json:"is_active"`
	ItemPackageDetailId int     `json:"item_package_detail_id"`
	ItemPackageId       int     `json:"item_package_id"`
	ItemId              int     `json:"item_id"`
	ItemClassId         int     `json:"item_class_id"`
	Quantity            float64 `json:"quantity"`
}

type UpdateitemPackageDetail struct {
	ItemPackageDetailId int     `json:"item_package_detail_id"`
	Quantity            float64 `json:"quantity"`
}

type ItemPackageDetailPayload struct {
	GetAllItemPackageResponse GetAllItemPackageResponse
	ItemPackageDetailResponse []ItemPackageDetailResponse
}
