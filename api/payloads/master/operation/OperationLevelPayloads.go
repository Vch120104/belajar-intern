package masteroperationpayloads

type OperationLevelGetAll struct {
	IsActive             bool   `json:"is_active"`
	OperationGroupCode   string `json:"operation_group_code"`
	OperationroupName    string `json:"operation_group_name"`
	OperationSectionCode string `json:"operation_section_code"`
	OperationSectionName string `json:"operation_section_name"`
	OperationKeyCode     string `json:"operation_key_code"`
	OperationKeyName     string `json:"operation_key_name"`
	OperationEntriesCode string `json:"operation_entries_code"`
	OperationEntriesName string `json:"operation_entries_name"`
	OperationEntriesId   int    `json:"operation_entries_id"`
}
