package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	generalservicepayloads "after-sales/api/payloads/cross-service/general-service"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type VatCompany struct {
	CompanyId          int    `json:"company_id"`
	CompanyPkpNumber   string `json:"pkp_no"`
	CompanyPkpDate     string `json:"pkp_date"`
	TaxName            string `json:"name"`
	CompanyTaxAddress  string `json:"address_street_1"`
	CompanyTaxAddress2 string `json:"address_street_2"`
	CompanyTaxAddress3 string `json:"address_street_3"`
	VillageId          int    `json:"village_id"`
	NpwpNumber         string `json:"npwp_no"`
	NpwpDate           string `json:"npwp_date"`
	TaxTransactionId   int    `json:"tax_transaction_id"`
	TaxBranchCode      string `json:"tax_branch_code"`
	TaxOfficeId        int    `json:"tax_office_id"`
}

type CompanyMasterDetailResponse struct {
	CompanyId           int        `json:"company_id"`
	CompanyCode         string     `json:"company_code"`
	CompanyName         string     `json:"company_name"`
	HeadOfficeCompanyId int        `json:"head_office_company_id"`
	VatCompany          VatCompany `json:"vat_company"`
}

type CompanyMasterResponse struct {
	CompanyId   int    `json:"company_id"`
	CompanyCode string `json:"company_code"`
	CompanyName string `json:"company_name"`
	IsDistbutor bool   `json:"is_distributor"`
	BizCategory string `json:"biz_category"`
}

func GetCompanyVat(id int) (VatCompany, *exceptions.BaseErrorResponse) {
	var getCompanyMaster VatCompany
	url := config.EnvConfigs.GeneralServiceUrl + "company-vat-by-id/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &getCompanyMaster)
	if err != nil {
		return getCompanyMaster, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error consuming external API for company VAT data",
			Err:        errors.New("error consuming external API for company VAT data"),
		}
	}
	return getCompanyMaster, nil
}

func GetCompanyDataById(companyId int) (generalservicepayloads.GetCompanyByIdResponses, *exceptions.BaseErrorResponse) {
	var companyResponse generalservicepayloads.GetCompanyByIdResponses
	companyUrl := config.EnvConfigs.GeneralServiceUrl + "company/" + strconv.Itoa(companyId)

	err := utils.CallAPI("GET", companyUrl, nil, &companyResponse)
	if err != nil {
		return companyResponse, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to fetch company data",
			Err:        errors.New("failed to fetch company data"),
		}
	}

	return companyResponse, nil
}

func GetCompanyByMultiId(ids []int, response interface{}) *exceptions.BaseErrorResponse {

	ids = utils.RemoveDuplicateIds(ids)
	validIds := make([]string, 0, len(ids))

	for _, id := range ids {
		if id != 0 {
			validIds = append(validIds, strconv.Itoa(id))
		}
	}

	strIds := "[" + strings.Join(validIds, ",") + "]"
	url := config.EnvConfigs.GeneralServiceUrl + "company-by-multi-id/" + strIds

	err := utils.CallAPI("GET", url, nil, response)
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error consuming external API for multiple company data",
			Err:        errors.New("error consuming external API for multiple company data"),
		}
	}
	return nil
}

func IsFTZCompany(companyId int) bool {
	return companyId == 139 //  ID for FTZ company check
}
