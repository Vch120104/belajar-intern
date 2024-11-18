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

type ServiceProfitCenter struct {
	ServiceProfitCenterId   int    `json:"service_profit_center_id"`
	ServiceProfitCenterCode string `json:"service_profit_center_code"`
	ServiceProfitCenterName string `json:"service_profit_center_description"`
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
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve service request status due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "service request status service is temporarily unavailable"
		}

		return getServiceRequestStatus, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting service request status by ID"),
		}
	}
	return getServiceRequestStatus, nil
}

func GetReferenceTypeById(id int) (ReferenceType, *exceptions.BaseErrorResponse) {
	var getReferenceType ReferenceType
	url := config.EnvConfigs.GeneralServiceUrl + "service-request-reference-type/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getReferenceType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve reference type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "reference type service is temporarily unavailable"
		}

		return getReferenceType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting reference type by ID"),
		}
	}
	return getReferenceType, nil
}

func GetServiceProfitCenterById(id int) (ServiceProfitCenter, *exceptions.BaseErrorResponse) {
	var getServiceProfitCenter ServiceProfitCenter
	url := config.EnvConfigs.GeneralServiceUrl + "service-profit-center/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getServiceProfitCenter)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve service profit center due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "service profit center service is temporarily unavailable"
		}

		return getServiceProfitCenter, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting service profit center by ID"),
		}
	}
	return getServiceProfitCenter, nil
}
