package masteritempayloads

type MarkupMasterResponse struct {
	IsActive          bool   `json:"is_active"`
	MarkupMasterId    int    `json:"markup_master_id"`
	MarkupCode        string `json:"markup_code"`
	MarkupDescription string `json:"markup_description"`
}

type MarkupMasterDropDownResponse struct {
	IsActive              bool   `json:"is_active"`
	MarkupMasterId        int    `json:"markup_master_id"`
	MarkupCodeDescription string `json:"markup_code_description"`
}
