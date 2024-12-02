package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"fmt"
	"net/http"
	"strconv"
)

type ItemClaimTypeMasterResponse struct {
	ItemClaimTypeDescription string `json:"item_claim_type_description"`
	ItemClaimTypeCode        string `json:"item_claim_type_code"`
	IsActive                 bool   `json:"is_active"`
	ItemClaimTypeId          int    `json:"item_claim_type_id"`
}

func GetItemClaimTypeMasterById(id int) (ItemClaimTypeMasterResponse, *exceptions.BaseErrorResponse) {
	var response ItemClaimTypeMasterResponse
	urlGetClaimTypeMaster := config.EnvConfigs.GeneralServiceUrl + "item-claim-type/" + strconv.Itoa(id)
	fmt.Println(urlGetClaimTypeMaster)
	err := utils.CallAPI("GET", urlGetClaimTypeMaster, nil, &response)
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get claim type master by id please check log",
			Err:        err,
		}
	}
	return response, nil
}
