package masteroperationpayloads

type OperationEntriesResponse struct {
	IsActive             bool   `json:"is_active" parent_entity:"mtr_operation_entries"`
	OperationEntriesId   int32  `json:"operation_entries_id" parent_entity:"mtr_operation_entries" main_table:"mtr_operation_entries"`
	OperationEntriesCode string `json:"operation_entries_code" parent_entity:"mtr_operation_entries"`
	OperationGroupId     int32  `json:"operation_group_id" parent_entity:"mtr_operation_group" references:"mtr_operation_group"`
	OperationSectionId   int32  `json:"operation_section_id" parent_entity:"mtr_operation_section" references:"mtr_operation_section"`
	OperationKeyId       int32  `json:"operation_key_id" parent_entity:"mtr_operation_key" references:"mtr_operation_key"`
	OperationEntriesDesc string `json:"operation_entries_desc" parent_entity:"mtr_operation_entries"`
}

type OperationEntriesRequest struct {
	OperationEntriesCode string `json:"operation_entries_code" parent_entity:"mtr_operation_entries"`
	OperationEntriesDesc string `json:"operation_entries_desc" parent_entity:"mtr_operation_entries"`
	OperationKeyId       int    `json:"operation_key_id" parent_entity:"mtr_operation_key" main_table:"mtr_operation_key"`
	OperationGroupId     int    `json:"operation_group_id" parent_entity:"mtr_operation_group" references:"mtr_operation_group"`
	OperationSectionId   int    `json:"operation_section_id" parent_entity:"mtr_operation_section" references:"mtr_operation_section"`
}
