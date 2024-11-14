package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type SalesOrderTransactionType struct {
	SoTransactionTypeId   int    `json:"sales_order_transaction_type_id"`
	SoTransactionTypeCode string `json:"sales_order_transaction_type_code"`
	SoTransactionTypeName string `json:"sales_order_transaction_type_name"`
}

func GetSoTransactionTypeByID(id int) (SalesOrderTransactionType, *exceptions.BaseErrorResponse) {
	var sotransactionType SalesOrderTransactionType
	url := config.EnvConfigs.GeneralServiceUrl + "sales-order-transaction-type/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &sotransactionType)
	if err != nil {
		return sotransactionType, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to retrieve transaction type",
			Err:        errors.New("error consuming external API while getting transaction type by ID"),
		}
	}
	return sotransactionType, nil
}
