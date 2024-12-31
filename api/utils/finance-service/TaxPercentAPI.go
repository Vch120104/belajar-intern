package financeserviceapiutils

import (
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type TaxPercentResponse struct {
	TaxPercent float64 `json:"tax_percent"`
}

const TaxPercentBaseUrl = "tax-fare/"

func GetTaxPercent(TaxServiceCode string, TaxTypeCode string, EffectiveDate time.Time) (TaxPercentResponse, *exceptions.BaseErrorResponse) {
	formattedTime := EffectiveDate.Format(time.RFC3339)
	getTaxPercentUrl := TaxPercentBaseUrl + "detail/tax-percent" +
		fmt.Sprintf("?tax_service_code=%s&tax_type_code=%s&effective_date=%s", TaxServiceCode, TaxTypeCode, formattedTime)

	TaxPercentResponse := TaxPercentResponse{}

	err := utils.CallAPI("GET", getTaxPercentUrl, nil, &TaxPercentResponse)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve tax percent due to an external service error"
		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "event service is temporarily unavailable"
		}
		return TaxPercentResponse, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting tax percent by tax type code"),
		}
	}
	return TaxPercentResponse, nil
}
