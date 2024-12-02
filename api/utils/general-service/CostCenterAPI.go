package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"net/http"
	"strconv"
)

type CostCenterResponses struct {
	CostCenterId   int    `json:"cost_center_id"`
	CostCenterCode string `json:"cost_center_code"`
	CostCenterName string `json:"cost_center_name"`
}

func GetCostCenterById(costCenterId int) (CostCenterResponses, *exceptions.BaseErrorResponse) {
	var GetCostCenter CostCenterResponses
	CostCenterURL := config.EnvConfigs.GeneralServiceUrl + "cost-center/" + strconv.Itoa(costCenterId)
	if err := utils.CallAPI("GET", CostCenterURL, nil, &GetCostCenter); err != nil {
		return GetCostCenter, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch Cost Center data from external service",
			Err:        err,
		}
	}
	return GetCostCenter, nil
}
