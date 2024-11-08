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
		return getItemRegulation, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching item regulation by id",
			Err:        errors.New("failed to retrieve item regulation data from external API by id"),
		}
	}
	return getItemRegulation, nil
}
