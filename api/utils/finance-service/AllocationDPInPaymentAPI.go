package financeserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
	"time"
)

type AllocationDpInPaymentResponse struct {
	AllocationDownPaymentInSystemNumber     int     `json:"allocation_down_payment_in_system_number" dataframe:"allocation_down_payment_in_system_number"`
	AllocationDownPaymentInDocumentNumber   string  `json:"allocation_down_payment_in_document_number" dataframe:"allocation_down_payment_in_document_number"`
	AllocationDownPaymentInDate             string  `json:"allocation_down_payment_in_date" dataframe:"allocation_down_payment_in_date"`
	AllocationCustomerId                    int     `json:"allocation_customer_id" dataframe:"allocation_customer_id"`
	AllocationDownPaymentInApprovalStatusId int     `json:"allocation_down_payment_in_approval_status_id" dataframe:"allocation_down_payment_in_approval_status_id"`
	DownPaymentInDocumentNumber             string  `json:"down_payment_in_document_number" dataframe:"down_payment_in_document_number"`
	CurrencyCode                            string  `json:"currency_code" dataframe:"currency_code"`
	TotalAmount                             float64 `json:"total_amount" dataframe:"total_amount"`
}

type CreateAllocationDPInRequest struct {
	CompanyId                   int       `json:"company_id" validate:"required"`
	AllocationDownPaymentInDate time.Time `json:"allocation_down_payment_in_date" validate:"required"`
	BrandId                     int       `json:"brand_id" validate:"required"`
	ProfitCenterId              int       `json:"profit_center_id" validate:"required"`
	DownPaymentInSystemNumber   int       `json:"down_payment_in_system_number" validate:"required"`
	AllocationCustomerId        int       `json:"allocation_customer_id" validate:"required"`
	Remark                      string    `json:"remark"`
}

func GetAllocationDpInPaymentById(id int) (AllocationDpInPaymentResponse, *exceptions.BaseErrorResponse) {
	var allocationDpInPayment AllocationDpInPaymentResponse
	url := config.EnvConfigs.FinanceServiceUrl + "allocation-down-payment-in/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &allocationDpInPayment)
	if err != nil {
		status := http.StatusBadGateway
		message := "Failed to retrieve allocation DP in payment due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "Allocation DP in payment service is temporarily unavailable"
		}

		return allocationDpInPayment, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting allocation DP in payment by ID"),
		}
	}
	return allocationDpInPayment, nil
}

func PostAllocationDpInPayment(request CreateAllocationDPInRequest) (*AllocationDpInPaymentResponse, *exceptions.BaseErrorResponse) {
	var response AllocationDpInPaymentResponse
	url := config.EnvConfigs.FinanceServiceUrl + "allocation-down-payment-in"
	err := utils.CallAPI("POST", url, request, &response)
	if err != nil {
		status := http.StatusBadGateway
		message := "Failed to create Allocation DP In Payment due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "Allocation DP In Payment service is temporarily unavailable"
		}

		return nil, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while creating Allocation DP In Payment"),
		}
	}
	return &response, nil
}
