package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"fmt"
	"net/http"
	"strconv"
)

type DealerRepresentativeResponse struct {
	DealerRepresentativeId   int    `json:"dealer_representative_id"`
	DealerRepresentativeCode int    `json:"dealer_representative_code"`
	DealerRepresentativeName string `json:"dealer_representative_name"`
}

type DealerRepresentativesResponse struct {
	DealerRepresentativeId   int    `json:"dealer_representative_id"`
	DealerRepresentativeCode string `json:"dealer_representative_code"`
	DealerRepresentativeName string `json:"dealer_representative_name"`
}

func GetDealerRepresentativeByMultiId(ids []int, abstractType interface{}) *exceptions.BaseErrorResponse {
	var strIds string = ""

	ids = utils.RemoveDuplicateIds(ids)
	for i := 0; i < len(ids); i++ {
		if ids[i] != 0 {
			strIds = strIds + strconv.Itoa(ids[i]) + ","
		}
	}
	if strIds != "" {
		strIds = "[" + strIds[:len(strIds)-1] + "]"
	} else {
		strIds = "[]"
	}
	url := config.EnvConfigs.GeneralServiceUrl + "dealer-representative-by-multi-id/" + strIds
	err := utils.CallAPI("GET", url, nil, &abstractType)
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed to get dealer representative by multiple id : %w", err),
		}
	}
	return nil
}

func GetDealerRepresentativeById(id int) (DealerRepresentativesResponse, *exceptions.BaseErrorResponse) {
	var dealer DealerRepresentativesResponse
	url := config.EnvConfigs.GeneralServiceUrl + "dealer-representative/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &dealer)
	if err != nil {
		return dealer, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed to get dealer representative by id : %w", err),
		}
	}
	return dealer, nil
}
