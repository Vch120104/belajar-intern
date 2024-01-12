package masteroperationpayloads

type OperationSectionResponse struct {
	IsActive                    bool   `json:"is_active"`
	OperationSectionId          int    `json:"operation_section_id"`
	OperationSectionCode        string `json:"operation_section_code"`
	OperationGroupId            int    `json:"operation_group_id"`
	OperationSectionDescription string `json:"operation_section_description"`
}

type OperationSectionListResponse struct {
	IsActive                    bool   `json:"is_active" parent_entity:"mtr_operation_section"`
	OperationSectionId          int    `json:"operation_section_id" parent_entity:"mtr_operation_section" main_table:"mtr_operation_section"`
	OperationSectionCode        string `json:"operation_section_code" parent_entity:"mtr_operation_section"`
	OperationSectionDescription string `json:"operation_section_description" parent_entity:"mtr_operation_section"`
	OperationGroupId            int    `json:"operation_group_id" parent_entity:"mtr_operation_group" references:"mtr_operation_group"`
	OperationGroupCode          string `json:"operation_group_code" parent_entity:"mtr_operation_group"`
	OperationGroupDescription   string `json:"operation_group_description" parent_entity:"mtr_operation_group"`
}

type OperationSectionCodeResponse struct {
	OperationSectionId   int    `json:"operation_section_id" parent_entity:"mtr_operation_section" main_table:"mtr_operation_section"`
	OperationGroupId     int    `json:"operation_group_id" parent_entity:"mtr_operation_group" references:"mtr_operation_group"`
	OperationSectionCode string `json:"operation_section_code" parent_entity:"mtr_operation_section"`
}
type OperationSectionNameResponse struct {
	OperationSectionId          int    `json:"operation_section_id" parent_entity:"mtr_operation_section" main_table:"mtr_operation_section"`
	OperationGroupId            int    `json:"operation_group_id" parent_entity:"mtr_operation_group" references:"mtr_operation_group"`
	OperationSectionDescription string `json:"operation_section_description" parent_entity:"mtr_operation_section"`
}
type OperationSectionRequest struct {
	IsActive                    bool   `json:"is_active"`
	OperationSectionId          int    `json:"operation_section_id"`
	OperationSectionCode        string `json:"operation_section_code"`
	OperationGroupId            int    `json:"operation_group_id"`
	OperationSectionDescription string `json:"operation_section_description"`
}
