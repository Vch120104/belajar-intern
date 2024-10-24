package generalservicepayloads

type SourceDocumentTypeMasterResponse struct {
	SourceDocumentTypeId          int    `json:"source_document_type_id"`
	SourceDocumentTypeCode        string `json:"source_document_type_code"`
	IsActive                      bool   `json:"is_active"`
	SourceDocumentTypeDescription string `json:"source_document_type_description"`
}
