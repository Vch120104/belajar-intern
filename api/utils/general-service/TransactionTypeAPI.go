package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"net/http"
	"strconv"
)

type TransactionTypeResponse struct {
	TransactionTypeName string `json:"transaction_type_name"`
	TransactionTypeCode string `json:"transaction_type_code"`
	TransactionTypeId   int    `json:"transaction_type_id"`
	IsActive            bool   `json:"is_active"`
}

func GetTransactionTypeById(id int) (TransactionTypeResponse, *exceptions.BaseErrorResponse) {
	var response TransactionTypeResponse

	url := config.EnvConfigs.GeneralServiceUrl + "transaction-type/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &response)
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get transaction type from general",
			Err:        err,
		}
	}
	return response, nil
}
