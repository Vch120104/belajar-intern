package userservice

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"time"
)

type ApprovalRequest struct {
	CompanyId         int       `json:"company_id"`
	ModuleId          int       `json:"module_id"`
	DocumentTypeId    int       `json:"document_type_id"`
	SourceDocNo       int       `json:"source_doc_no"`
	SourceSysNo       int       `json:"source_sys_no"`
	SourceAmount      float64   `json:"source_amount"`
	SourceDate        time.Time `json:"source_date"`
	TransactionTypeId int       `json:"transaction_type_id"`
	BrandId           int       `json:"brand_id"`
	ProfitCenterId    int       `json:"profit_center_id"`
	CostCenterId      int       `json:"cost_center_id"`
	RequestBy         int       `json:"request_by"`
	IsVoid            bool      `json:"is_void"`
}

type ApprovalResponse struct {
	ApprovalRequestId int `json:"approval_request_id"`
}

func CreateApproval(request ApprovalRequest) (ApprovalResponse, *exceptions.BaseErrorResponse) {
	var response ApprovalResponse
	url := config.EnvConfigs.UserServiceUrl + "approval-request"
	err := utils.CallAPI("POST", url, request, &response)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve brand due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "brand service is temporarily unavailable"
		}

		return response, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting brand by ID"),
		}
	}
	return response, nil
}
