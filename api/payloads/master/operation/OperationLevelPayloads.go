package masteroperationpayloads

// type OperationLevelGetAll struct {
// 	IsActive                    bool   `json:"is_active"`
// 	OperationLevelId            int    `json:"operation_level_id"`
// 	OperationGroupCode          string `json:"operation_group_code"`
// 	OperationGroupDescription   string `json:"operation_group_description"`
// 	OperationSectionCode        string `json:"operation_section_code"`
// 	OperationSectionDescription string `json:"operation_section_description"`
// 	OperationKeyCode            string `json:"operation_key_code"`
// 	OperationKeyDescription     string `json:"operation_key_description"`
// 	OperationEntriesCode        string `json:"operation_entries_code"`
// 	OperationEntriesDescription string `json:"operation_entries_description"`
// 	OperationEntriesId          int    `json:"operation_entries_id"`
// }

type OperationLevelGetAll struct {
	IsActive                    bool   `json:"is_active" parent_entity:"mtr_operation_level" main_table:"mtr_operation_level"`
	OperationLevelId            int    `json:"operation_level_id" parent_entity:"mtr_operation_level" main_table:"mtr_operation_level"`
	OperationEntriesId          int    `json:"operation_entries_id" parent_entity:"mtr_operation_entries" references:"mtr_operation_entries"`
	OperationEntriesCode        string `json:"operation_entries_code" parent_entity:"mtr_operation_entries"`
	OperationEntriesDescription string `json:"operation_entries_description" parent_entity:"mtr_operation_entries"`
	OperationGroupId            int    `json:"operation_group_id" parent_entity:"mtr_operation_group" references:"mtr_operation_group"`
	OperationGroupCode          string `json:"operation_group_code" parent_entity:"mtr_operation_group"`
	OperationGroupDescription   string `json:"operation_group_description" parent_entity:"mtr_operation_group"`
	OperationSectionId          int    `json:"operation_section_id" parent_entity:"mtr_operation_section" references:"mtr_operation_section"`
	OperationSectionCode        string `json:"operation_section_code" parent_entity:"mtr_operation_section"`
	OperationSectionDescription string `json:"operation_section_description" parent_entity:"mtr_operation_section"`
	OperationKeyId              int    `json:"operation_key_id" parent_entity:"mtr_operation_key" references:"mtr_operation_key"`
	OperationKeyCode            string `json:"operation_key_code" parent_entity:"mtr_operation_key"`
	OperationKeyDescription     string `json:"operation_key_description" parent_entity:"mtr_operation_key"`
}

type OperationLevelRequest struct {
	OperationModelMappingId int  `json:"operation_model_mapping_id"`
	OperationEntriesId      int  `json:"operation_entries_id"`
	IsActive                bool `json:"is_active"`
	OperationLevelId        int  `json:"operation_level_id"`
	// OperationEntriesCode    string `json:"operation_entries_code"`
	// OperationKeyCode        string `json:"operation_key_code"`
	// OperationGroupCode      string `json:"operation_group_code"`
	// OperationSectionCode    string `json:"operation_section_code"`
}
