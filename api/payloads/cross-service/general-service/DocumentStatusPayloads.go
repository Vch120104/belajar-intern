package generalservicepayloads

type DocumentStatusPayloads struct {
	DocumentStatusCode        string `json:"document_status_code"`
	DocumentStatusDescription string `json:"document_status_description"`
	IsActive                  bool   `json:"is_active"`
	DocumentStatusId          int    `json:"document_status_id"`
}
