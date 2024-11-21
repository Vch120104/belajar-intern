package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type WarrantyClaimTypeResponse struct {
	IsActive                     bool   `json:"is_active"`
	WarrantyClaimTypeId          int    `json:"warranty_claim_type_id"`
	WarrantyClaimTypeCode        string `json:"warranty_claim_type_code"`
	WarrantyClaimTypeDescription string `json:"warranty_claim_type_description"`
}

func GetWarrantyClaimTypeById(id int) (WarrantyClaimTypeResponse, *exceptions.BaseErrorResponse) {
	var getWarrantyClaimType WarrantyClaimTypeResponse
	url := config.EnvConfigs.GeneralServiceUrl + "warranty-claim-type/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getWarrantyClaimType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve warranty claim type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "warranty claim type service is temporarily unavailable"
		}

		return getWarrantyClaimType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting warranty claim type by ID"),
		}
	}

	return getWarrantyClaimType, nil
}
