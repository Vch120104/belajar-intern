package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type CompanyReferenceBetByIdResponse struct {
	CurrencyId                int         `json:"currency_id"`
	CoaGroupId                int         `json:"coa_group_id"`
	OperationDiscountOuterKpp json.Number `json:"operation_discount_outer_kpp"`
	MarginOuterKpp            json.Number `json:"margin_outer_kpp"`
	AdjustmentReasonId        int         `json:"adjustment_reason_id"`
	LeadTimeUnitEtd           int         `json:"lead_time_unit_etd"`
	BankAccReceiveCompanyId   int         `json:"bank_acc_receive_company_id"`
	UnitWarehouseId           int         `json:"unit_warehouse_id"`
	TimeDifference            int         `json:"time_difference"`
	UseDms                    bool        `json:"use_dms"`
	UseJpcb                   bool        `json:"use_jpcb"`
	CheckMonthEnd             bool        `json:"check_month_end"`
	IsDistributor             bool        `json:"is_distributor"`
	WithVat                   bool        `json:"with_vat"`
}

func GetCompanyReferenceById(id int) (CompanyReferenceBetByIdResponse, *exceptions.BaseErrorResponse) {
	CompanyReferenceBetByIdResponseData := CompanyReferenceBetByIdResponse{}

	CompanyReferenceUrl := fmt.Sprintf("%scompany-reference/%s", config.EnvConfigs.GeneralServiceUrl, strconv.Itoa(id))
	errFetchCompany := utils.CallAPI("GET", CompanyReferenceUrl, nil, &CompanyReferenceBetByIdResponseData)
	if errFetchCompany != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve company references due to an external service error"

		if errors.Is(errFetchCompany, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "company references service is temporarily unavailable"
		}

		return CompanyReferenceBetByIdResponseData, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting company references by ID"),
		}
	}
	return CompanyReferenceBetByIdResponseData, nil
}
