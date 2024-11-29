package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type WarrantyFreeServiceTypeResponse struct {
	WarrantyFreeServiceTypeId   int    `json:"warranty_free_service_type_id"`
	WarrantyFreeServiceTypeCode string `json:"warranty_free_service_type_code"`
	WarrantyFreeServiceTypeName string `json:"warranty_free_service_type_description"`
}

func GetWarrantyFreeServiceTypeById(id int) (WarrantyFreeServiceTypeResponse, *exceptions.BaseErrorResponse) {
	var wfst WarrantyFreeServiceTypeResponse
	url := config.EnvConfigs.GeneralServiceUrl + "warranty-free-service-type/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &wfst)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve warranty free service types due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "warranty free service types service is temporarily unavailable"
		}

		return wfst, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting warranty free service types by ID"),
		}
	}
	return wfst, nil
}
