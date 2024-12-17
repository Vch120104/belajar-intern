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

func GetSoTransactionTypeById(id int) (SalesOrderTransactionType, *exceptions.BaseErrorResponse) {
	var sotransactionType SalesOrderTransactionType
	url := config.EnvConfigs.GeneralServiceUrl + "sales-order-transaction-type/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &sotransactionType)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve sales order transaction type due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "sales order transaction type service is temporarily unavailable"
		}

		return sotransactionType, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting sales order transaction type by ID"),
		}
	}
	return sotransactionType, nil
}
