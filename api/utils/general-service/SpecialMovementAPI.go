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
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve special movement due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "special movement service is temporarily unavailable"
		}

		return getSpecialMovement, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting special movement by ID"),
		}
	}
	return getSpecialMovement, nil
}
