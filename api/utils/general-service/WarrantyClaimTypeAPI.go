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
		return getWarrantyClaimType, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching warranty claim type by id",
			Err:        errors.New("failed to retrieve warranty claim type data from external API by id"),
		}
	}

	return getWarrantyClaimType, nil
}
