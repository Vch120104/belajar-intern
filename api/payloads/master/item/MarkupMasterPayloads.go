package masteritempayloads

type MarkupMasterResponse struct {
	IsActive                bool   `json:"is_active"`
	MarkupMasterId          int    `json:"markup_master_id"`
	MarkupMasterCode        string `json:"markup_master_code"`
	MarkupMasterDescription string `json:"markup_master_description"`
}

type MarkupMasterDropDownResponse struct {
	IsActive                bool   `json:"is_active"`
	MarkupMasterId          int    `json:"markup_master_id"`
	MarkupMasterCodeDescription        string `json:"markup_master_code_description"`
}