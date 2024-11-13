package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type DocumentStatusPayloads struct {
	DocumentStatusCode        string `json:"document_status_code"`
	DocumentStatusDescription string `json:"document_status_description"`
	IsActive                  bool   `json:"is_active"`
	DocumentStatusId          int    `json:"document_status_id"`
}

func GetDocumentStatusById(id int) (DocumentStatusPayloads, *exceptions.BaseErrorResponse) {
	var DocResponse DocumentStatusPayloads
	DocumentStatusUrl := config.EnvConfigs.GeneralServiceUrl + "document-status/" + strconv.Itoa(id)
	//err := utils.Get(DocumentStatusUrl, &DocResponse, nil)
	err := utils.CallAPI("GET", DocumentStatusUrl, nil, &DocResponse)
	if err != nil {
		return DocResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed To fetch Document Status From External Service",
			Err:        errors.New("failed To fetch Document Status From External Service"),
		}
	}
	return DocResponse, nil
}
func GetDocumentStatusByCode(code string) (DocumentStatusPayloads, *exceptions.BaseErrorResponse) {
	var DocResponse DocumentStatusPayloads
	DocumentStatusUrl := config.EnvConfigs.GeneralServiceUrl + "document-status-by-code/" + code
	//err := utils.Get(DocumentStatusUrl, &DocResponse, nil)
	err := utils.CallAPI("GET", DocumentStatusUrl, nil, &DocResponse)
	if err != nil {
		return DocResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed To fetch Document Status From External Service",
			Err:        errors.New("failed To fetch Document Status From External Service"),
		}
	}
	return DocResponse, nil
}

type SourceDocumentTypeMasterResponse struct {
	SourceDocumentTypeId          int    `json:"source_document_type_id"`
	SourceDocumentTypeCode        string `json:"source_document_type_code"`
	IsActive                      bool   `json:"is_active"`
	SourceDocumentTypeDescription string `json:"source_document_type_description"`
}

func GetDocumentTypeByCode(code string) (SourceDocumentTypeMasterResponse, *exceptions.BaseErrorResponse) {
	var SourceDocType SourceDocumentTypeMasterResponse
	SourceDocTypeUrl := config.EnvConfigs.GeneralServiceUrl + "source-document-type-code/" + code
	err := utils.CallAPI("GET", SourceDocTypeUrl, nil, &SourceDocType)
	//if err := utils.Get(SourceDocTypeUrl, &SourceDocType, nil)
	if err != nil {
		return SourceDocType, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("failed To fetch source document Status From External Service"),
		}
	}
	return SourceDocType, nil
}
