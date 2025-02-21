package financeserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// TaxPercentParams untuk query parameter API tax percent
type TaxPercentParams struct {
	TaxServiceCode string    `json:"tax_service_code"`
	TaxTypeCode    string    `json:"tax_type_code"`
	EffectiveDate  time.Time `json:"effective_date"`
}

// TaxPercentResponse untuk respons API tax percent
type TaxPercentResponse struct {
	TaxPercent float64 `json:"tax_percent"`
}

// GetTaxPercent mengambil nilai pajak berdasarkan kode layanan pajak, tipe pajak, dan tanggal efektif
func GetTaxPercent(params TaxPercentParams) (TaxPercentResponse, *exceptions.BaseErrorResponse) {
	var response TaxPercentResponse

	baseURL := config.EnvConfigs.FinanceServiceUrl + "tax-fare/detail/tax-percent"

	formattedDate := params.EffectiveDate.Format(time.RFC3339)
	queryParams := fmt.Sprintf("?tax_service_code=%s&tax_type_code=%s&effective_date=%s",
		params.TaxServiceCode, params.TaxTypeCode, formattedDate)

	url := baseURL + queryParams

	err := utils.CallAPI("GET", url, nil, &response)
	if err != nil {
		status := http.StatusBadGateway // Default error status 502
		message := "Failed to retrieve tax percent due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "tax service is temporarily unavailable"
		}

		return response, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting tax percent"),
		}
	}

	return response, nil
}
