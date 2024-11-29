package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type TermOfPaymentResponse struct {
	TermOfPaymentId   int    `json:"term_of_payment_id"`
	TermOfPaymentCode string `json:"term_of_payment_code"`
	TermOfPaymentName string `json:"term_of_payment_name"`
}

func GetTermOfPaymentById(id int) (TermOfPaymentResponse, *exceptions.BaseErrorResponse) {
	var top TermOfPaymentResponse
	url := config.EnvConfigs.GeneralServiceUrl + "term-of-payment/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &top)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve top due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "top service is temporarily unavailable"
		}

		return top, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting top by ID"),
		}
	}
	return top, nil
}

func GetTermOfPaymentByCode(code string) (TermOfPaymentResponse, *exceptions.BaseErrorResponse) {
	var top TermOfPaymentResponse
	url := config.EnvConfigs.GeneralServiceUrl + "term-of-payment-code/" + code
	err := utils.CallAPI("GET", url, nil, &top)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve top due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "top service is temporarily unavailable"
		}

		return top, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting top by ID"),
		}
	}
	return top, nil
}
