package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
	"time"
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

type DocumentMasterRequest struct {
	CompanyId         int       `json:"company_id"`
	TransactionDate   time.Time `json:"transaction_date"`
	DocumentTypeId    int       `json:"document_type_id"`
	BrandId           int       `json:"brand_id"`
	ProfitCenterId    int       `json:"profit_center_id"`
	TransactionTypeId int       `json:"transaction_type_id"`
	BankCompanyId     int       `json:"bank_company_id"`
}

type DocumentMasterResponse struct {
	GeneratedDocumentNumber string `json:"generated_document_number"`
}

func GetDocumentNumber(request DocumentMasterRequest) (DocumentMasterResponse, *exceptions.BaseErrorResponse) {
	var response DocumentMasterResponse
	url := config.EnvConfigs.GeneralServiceUrl + "last-document-number?company_id=" + strconv.Itoa(request.CompanyId) +
		"&transaction_date=" + request.TransactionDate.Format(utils.RFC3339) + "&document_type_id=" + strconv.Itoa(request.DocumentTypeId)
	if request.ProfitCenterId != 0 {
		url = url + "&profit_center_id=" + strconv.Itoa(request.ProfitCenterId)
	}
	if request.BrandId != 0 {
		url = url + "&brand_id=" + strconv.Itoa(request.BrandId)
	}
	if request.BankCompanyId != 0 {
		url = url + "&bank_company_id=" + strconv.Itoa(request.BankCompanyId)
	}
	if request.TransactionTypeId != 0 {
		url = url + "&transaction_type_id=" + strconv.Itoa(request.TransactionTypeId)
	}

	err := utils.CallAPI("GET", url, nil, &response)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve source document status due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "source document status service is temporarily unavailable"
		}

		return response, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting source document status by ID"),
		}
	}
	return response, nil
}
