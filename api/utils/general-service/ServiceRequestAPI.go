package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type ServiceRequestStatus struct {
	ServiceRequestStatusID          int    `json:"service_request_reference_status_id"`
	ServiceRequestStatusCode        string `json:"service_request_reference_status_code"`
	ServiceRequestStatusDescription string `json:"service_request_reference_status_description"`
}

type ReferenceType struct {
	ReferenceTypeId   int    `json:"service_request_reference_type_id"`
	ReferenceTypeCode string `json:"service_request_reference_type_code"`
	ReferenceTypeName string `json:"service_request_reference_type_description"`
}

func GetServiceRequestStatusById(id int) (ServiceRequestStatus, *exceptions.BaseErrorResponse) {
	var getServiceRequestStatus ServiceRequestStatus
	url := config.EnvConfigs.GeneralServiceUrl + "service-request-status/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getServiceRequestStatus)
	if err != nil {
		return getServiceRequestStatus, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching service request status by id",
			Err:        errors.New("failed to retrieve service request status data from external API by id"),
		}
	}
	return getServiceRequestStatus, nil
}

func GetReferenceTypeById(id int) (ReferenceType, *exceptions.BaseErrorResponse) {
	var getReferenceType ReferenceType
	url := config.EnvConfigs.GeneralServiceUrl + "service-request-reference-type/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getReferenceType)
	if err != nil {
		return getReferenceType, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching reference type by id",
			Err:        errors.New("failed to retrieve reference type data from external API by id"),
		}
	}
	return getReferenceType, nil
}
