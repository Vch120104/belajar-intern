package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type ItemRegulationResponse struct {
	IsActive           bool   `json:"is_active"`
	ItemRegulationId   int    `json:"item_regulation_id"`
	ItemRegulationCode string `json:"item_regulation_code"`
	ItemRegulationName string `json:"item_regulation_nmae"`
}

func GetItemRegulationById(id int) (ItemRegulationResponse, *exceptions.BaseErrorResponse) {
	var getItemRegulation ItemRegulationResponse
	url := config.EnvConfigs.GeneralServiceUrl + "item-regulation/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getItemRegulation)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve item regulation due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "item regulation service is temporarily unavailable"
		}

		return getItemRegulation, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting item regulation by ID"),
		}
	}
	return getItemRegulation, nil
}
