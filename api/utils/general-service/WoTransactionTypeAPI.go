package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type WorkOrderTransactionType struct {
	WoTransactionTypeId   int    `json:"work_order_transaction_type_id"`
	WoTransactionTypeCode string `json:"work_order_transaction_type_code"`
	WoTransactionTypeName string `json:"work_order_transaction_type_name"`
}

func GetWoTransactionTypeById(id int) (WorkOrderTransactionType, *exceptions.BaseErrorResponse) {
	var transactionType WorkOrderTransactionType
	url := config.EnvConfigs.GeneralServiceUrl + "work-order-transaction-type/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &transactionType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve work order type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "Transaction type service is temporarily unavailable"
		}

		return transactionType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting transaction type by ID"),
		}
	}
	return transactionType, nil
}

func GetWoTransactionTypeByCode(code string) (WorkOrderTransactionType, *exceptions.BaseErrorResponse) {
	var transactionType WorkOrderTransactionType
	url := config.EnvConfigs.GeneralServiceUrl + "work-order-transaction-type-by-code/" + code

	err := utils.CallAPI("GET", url, nil, &transactionType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve work order type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "Transaction type service is temporarily unavailable"
		}

		return transactionType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting transaction type by code"),
		}
	}
	return transactionType, nil
}
