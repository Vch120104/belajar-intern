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

const (
	errorMsgCode    = "error consuming external API to get approval status by code"
	errorMsgID      = "error consuming external API to get approval status by ID"
	errorMsgMultiID = "error consuming external API to get approval status by multiple IDs"
)

func GetApprovalStatusByCode(code string) (ApprovalStatusResponse, *exceptions.BaseErrorResponse) {
	var getApprovalStatusTemp ApprovalStatusTempResponse
	var getApprovalStatus ApprovalStatusResponse
	url := config.EnvConfigs.GeneralServiceUrl + "approval-status-by-code/" + code

	err := utils.CallAPI("GET", url, nil, &getApprovalStatusTemp)
	if err != nil {
		return getApprovalStatus, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    errorMsgCode,
			Err:        errors.New(errorMsgCode),
		}
	}

	getApprovalStatus = ApprovalStatusResponse{
		ApprovalStatusId:          getApprovalStatusTemp.ApprovalStatusId,
		ApprovalStatusCode:        strconv.Itoa(getApprovalStatusTemp.ApprovalStatusCode),
		ApprovalStatusDescription: getApprovalStatusTemp.ApprovalStatusDescription,
	}
	return getApprovalStatus, nil
}

func GetApprovalStatusByID(id int) (ApprovalStatusResponse, *exceptions.BaseErrorResponse) {
	var getApprovalStatusTemp ApprovalStatusTempResponse
	var getApprovalStatus ApprovalStatusResponse
	url := config.EnvConfigs.GeneralServiceUrl + "approval-status/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getApprovalStatusTemp)
	if err != nil {
		return getApprovalStatus, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    errorMsgID,
			Err:        errors.New(errorMsgID),
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

	url := config.EnvConfigs.GeneralServiceUrl + "approval-status-by-multi-id/" + strIds
	err := utils.CallAPI("GET", url, nil, response)
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    errorMsgMultiID,
			Err:        errors.New(errorMsgMultiID),
		}
	}

	return nil
}
