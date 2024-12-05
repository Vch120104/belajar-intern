package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
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
	url := config.EnvConfigs.GeneralServiceUrl + "dealer-representative-multi-id/" + strIds
	err := utils.CallAPI("GET", url, nil, &abstractType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve dealer representative due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "dealer representative service is temporarily unavailable"
		}

		return &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting dealer representative by ID"),
		}
	}
	return nil
}

func GetDealerRepresentativeById(id int) (DealerRepresentativesResponse, *exceptions.BaseErrorResponse) {
	var dealer DealerRepresentativesResponse
	url := config.EnvConfigs.GeneralServiceUrl + "dealer-representative/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &dealer)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve dealer representative due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "dealer representative service is temporarily unavailable"
		}

		return dealer, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting dealer representative by ID"),
		}
	}
	return dealer, nil
}
