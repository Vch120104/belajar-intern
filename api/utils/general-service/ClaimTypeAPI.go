package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type ClaimTypeResponse struct {
	IsActive             bool   `json:"is_active"`
	ClaimTypeId          int    `json:"claim_type_id"`
	ClaimTypeCode        string `json:"claim_type_code"`
	ClaimTypeDescription string `json:"claim_type_description"`
}

func GetClaimTypeById(id int) (ClaimTypeResponse, *exceptions.BaseErrorResponse) {
	var getClaimType ClaimTypeResponse
	url := config.EnvConfigs.GeneralServiceUrl + "claim-type/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getClaimType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve claim type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "claim type service is temporarily unavailable"
		}

		return getClaimType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting claim type by ID"),
		}
	}
	return getClaimType, nil
}
