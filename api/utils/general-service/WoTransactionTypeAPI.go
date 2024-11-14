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

func GetTransactionTypeByID(id int) (WorkOrderTransactionType, *exceptions.BaseErrorResponse) {
	var transactionType WorkOrderTransactionType
	url := config.EnvConfigs.GeneralServiceUrl + "work-order-transaction-type/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &transactionType)
	if err != nil {
		return transactionType, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve transaction type",
			Err:        errors.New("error consuming external API while getting transaction type by ID"),
		}
	}
	return transactionType, nil
}

func GetTransactionTypeByCode(code string) (WorkOrderTransactionType, *exceptions.BaseErrorResponse) {
	var transactionType WorkOrderTransactionType
	url := config.EnvConfigs.GeneralServiceUrl + "work-order-transaction-type/by-code/" + code

	err := utils.CallAPI("GET", url, nil, &transactionType)
	if err != nil {
		return transactionType, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve transaction type",
			Err:        errors.New("error consuming external API while getting transaction type by code"),
		}
	}
	return transactionType, nil
}
