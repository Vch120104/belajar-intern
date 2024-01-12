package masteritempayloads

type ItemClassRequest struct {
	IsActive      bool   `json:"is_active"`
	ItemClassCode string `json:"item_class_code"`
	ItemGroupID   int    `json:"item_group_id"` //FK with mtr_item_group common-general service
	LineTypeID    int    `json:"line_type_id"`  //FK with mtr_line_type common-general service
	ItemClassName string `json:"item_class_name"`
}

type ItemClassResponse struct {
	IsActive      bool   `json:"is_active" parent_entity:"mtr_item_class"`
	ItemClassId   int    `json:"item_class_id" parent_entity:"mtr_item_class"`
	ItemClassCode string `json:"item_class_code" parent_entity:"mtr_item_class"`
	ItemGroupId   int    `json:"item_group_id"` //FK with mtr_item_group common-general service
	LineTypeId    int    `json:"line_type_id"`  //FK with mtr_line_type common-general service
	ItemClassName string `json:"item_class_name" parent_entity:"mtr_item_class"`
}

// IsActive      bool   `json:"is_active" parent_entity:"mtr_item"`
// ItemId        int    `json:"item_id" parent_entity:"mtr_item" main_table:"mtr_item"`
// ItemCode      string `json:"item_code" parent_entity:"mtr_item"`
// ItemName      string `json:"item_name" parent_entity:"mtr_item"`
// ItemType      string `json:"item_type" parent_entity:"mtr_item"`
// ItemGroupId   int    `json:"item_group_id" parent_entity:"mtr_item"`                                   //fk luar mtr_item_group -> item_group_name
// ItemClassId   int    `json:"item_class_id" parent_entity:"mtr_item_class" references:"mtr_item_class"` //fk dalam item_class_id -> ItemClassName
// ItemClassCode string `json:"item_class_code" parent_entity:"mtr_item_class"`
// SupplierId    int    `json:"supplier_id" parent_entity:"mtr_item"` //fk luar mtr_supplier, supplier_code dan supplier_name
