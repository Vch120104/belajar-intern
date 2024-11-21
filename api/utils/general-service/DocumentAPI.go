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
	err := utils.CallAPI("GET", DocumentStatusUrl, nil, &DocResponse)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve document status due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "document status service is temporarily unavailable"
		}

		return DocResponse, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting document status by ID"),
		}
	}
	return DocResponse, nil
}
func GetDocumentStatusByCode(code string) (DocumentStatusPayloads, *exceptions.BaseErrorResponse) {
	var DocResponse DocumentStatusPayloads
	DocumentStatusUrl := config.EnvConfigs.GeneralServiceUrl + "document-status-by-code/" + code
	err := utils.CallAPI("GET", DocumentStatusUrl, nil, &DocResponse)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve document status due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "document status service is temporarily unavailable"
		}

		return DocResponse, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting document status by ID"),
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
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve source document status due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "source document status service is temporarily unavailable"
		}

		return SourceDocType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting source document status by ID"),
		}
	}
	return SourceDocType, nil
}
