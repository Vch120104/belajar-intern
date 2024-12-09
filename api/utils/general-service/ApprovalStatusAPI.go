package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type ApprovalStatusTempResponse struct {
	ApprovalStatusId          int    `json:"approval_status_id"`
	ApprovalStatusCode        int    `json:"approval_status_code"`
	ApprovalStatusDescription string `json:"approval_status_description"`
}

type ApprovalStatusResponse struct {
	ApprovalStatusId          int    `json:"approval_status_id"`
	ApprovalStatusCode        string `json:"approval_status_code"`
	ApprovalStatusDescription string `json:"approval_status_description"`
}

func GetApprovalStatusByCode(code string) (ApprovalStatusResponse, *exceptions.BaseErrorResponse) {
	var getApprovalStatusTemp ApprovalStatusTempResponse
	var getApprovalStatus ApprovalStatusResponse
	url := config.EnvConfigs.GeneralServiceUrl + "approval-status-code/" + code

	err := utils.CallAPI("GET", url, nil, &getApprovalStatusTemp)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve approval status due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "approval status service is temporarily unavailable"
		}

		return getApprovalStatus, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting approval status by ID"),
		}
	}

	getApprovalStatus = ApprovalStatusResponse{
		ApprovalStatusId:          getApprovalStatusTemp.ApprovalStatusId,
		ApprovalStatusCode:        strconv.Itoa(getApprovalStatusTemp.ApprovalStatusCode),
		ApprovalStatusDescription: getApprovalStatusTemp.ApprovalStatusDescription,
	}
	return getApprovalStatus, nil
}

func GetApprovalStatusById(id int) (ApprovalStatusResponse, *exceptions.BaseErrorResponse) {
	var getApprovalStatusTemp ApprovalStatusTempResponse
	var getApprovalStatus ApprovalStatusResponse
	url := config.EnvConfigs.GeneralServiceUrl + "approval-status/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getApprovalStatusTemp)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve approval status due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "approval status service is temporarily unavailable"
		}

		return getApprovalStatus, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting approval status by ID"),
		}
	}

	getApprovalStatus = ApprovalStatusResponse{
		ApprovalStatusId:          getApprovalStatusTemp.ApprovalStatusId,
		ApprovalStatusCode:        strconv.Itoa(getApprovalStatusTemp.ApprovalStatusCode),
		ApprovalStatusDescription: getApprovalStatusTemp.ApprovalStatusDescription,
	}
	return getApprovalStatus, nil
}

func GetApprovalStatusByMultiId(ids []int, response interface{}) *exceptions.BaseErrorResponse {

	ids = utils.RemoveDuplicateIds(ids)

	var idStrings []string
	for _, id := range ids {
		if id != 0 {
			idStrings = append(idStrings, strconv.Itoa(id))
		}
	}
	strIds := "[" + strings.Join(idStrings, ",") + "]"

	url := config.EnvConfigs.GeneralServiceUrl + "approval-status-multi-id/" + strIds
	err := utils.CallAPI("GET", url, nil, response)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve approval status due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "approval status service is temporarily unavailable"
		}

		return &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting approval status by ID"),
		}
	}

	return nil
}
