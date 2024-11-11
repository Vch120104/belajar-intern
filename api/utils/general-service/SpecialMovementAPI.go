package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type SpecialMovementResponse struct {
	IsActive            bool   `json:"is_active"`
	SpecialMovementId   int    `json:"special_movement_id"`
	SpecialMovementCode string `json:"special_movement_code"`
	SpecialMovementName string `json:"special_movement_name"`
}

func GetSpecialMovementById(id int) (SpecialMovementResponse, *exceptions.BaseErrorResponse) {
	var getSpecialMovement SpecialMovementResponse
	url := config.EnvConfigs.GeneralServiceUrl + "special-movement/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getSpecialMovement)
	if err != nil {
		return getSpecialMovement, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching special movement by id",
			Err:        errors.New("failed to retrieve special movement data from external API by id"),
		}
	}
	return getSpecialMovement, nil
}
