package masteroperationpayloads

type OperationKeyResponse struct {
	IsActive                bool   `json:"is_active"`
	OperationKeyId          int    `json:"operation_key_id"`
	OperationKeyCode        string `json:"operation_key_code" validate:"required,max=5"`
	OperationGroupId        int    `json:"operation_group_id" validate:"required"`
	OperationSectionId      int    `json:"operation_section_id" validate:"required"`
	OperationKeyDescription string `json:"operation_key_description" validate:"required"`
}

type OperationkeyListResponse struct {
	OperationGroupId            int    `json:"operation_group_id" parent_entity:"mtr_operation_group" references:"mtr_operation_group"`
	OperationGroupCode          string `json:"operation_group_code" parent_entity:"mtr_operation_group"`
	OperationGroupDescription   string `json:"operation_group_description" parent_entity:"mtr_operation_group"`
	IsActive                    bool   `json:"is_active" parent_entity:"mtr_operation_key"`
	OperationSectionId          int    `json:"operation_section_id" parent_entity:"mtr_operation_section" references:"mtr_operation_section"`
	OperationSectionCode        string `json:"operation_section_code" parent_entity:"mtr_operation_section"`
	OperationSectionDescription string `json:"operation_section_description" parent_entity:"mtr_operation_section"`
	OperationKeyId              int    `json:"operation_key_id" parent_entity:"mtr_operation_key" main_table:"mtr_operation_key"`
	OperationKeyCode            string `json:"operation_key_code" parent_entity:"mtr_operation_key"`
	OperationKeyDescription     string `json:"operation_key_description" parent_entity:"mtr_operation_key" `
}

type OperationKeyRequest struct {
	OperationKeyCode        string `json:"operation_key_code"`
	OperationGroupId        int    `json:"operation_group_id"`
	OperationSectionId      int    `json:"operation_section_id"`
	OperationKeyDescription string `json:"operation_key_description"`
}

type OperationKeyNameResponse struct {
	OperationKeyId          int    `json:"operation_key_id" parent_entity:"mtr_operation_key" main_table:"mtr_operation_key"`
	OperationKeyCode        string `json:"operation_key_code" parent_entity:"mtr_operation_key"`
	OperationGroupId        int    `json:"operation_group_id" parent_entity:"mtr_operation_group" references:"mtr_operation_group"`
	OperationSectionId      int    `json:"operation_section_id" parent_entity:"mtr_operation_section" references:"mtr_operation_section"`
	OperationKeyDescription string `json:"operation_key_description" parent_entity:"mtr_operation_key"`
}

type OperationKeyCodeResponse struct {
	OperationKeyId          int    `json:"operation_key_id" parent_entity:"mtr_operation_key" main_table:"mtr_operation_key"`
	OperationKeyCode        string `json:"operation_key_code" parent_entity:"mtr_operation_key"`
	OperationGroupId        int    `json:"operation_group_id" parent_entity:"mtr_operation_group" references:"mtr_operation_group"`
	OperationSectionId      int    `json:"operation_section_id" parent_entity:"mtr_operation_section" references:"mtr_operation_section"`
	OperationKeyDescription string `json:"operation_key_description" parent_entity:"mtr_operation_key"`
}
