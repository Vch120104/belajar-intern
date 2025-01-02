package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"reflect"
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
type GetCompanyByIdResponses struct {
	CompanyName            string  `json:"company_name"`
	RegionId               int     `json:"region_id"`
	TermOfPaymentId        int     `json:"term_of_payment_id"`
	TaxCompanyId           int     `json:"tax_company_id"`
	FinanceAreaId          int     `json:"finance_area_id"`
	BusinessScopeId        int     `json:"business_scope_id"`
	IsActive               bool    `json:"is_active"`
	CompanyId              int     `json:"company_id"`
	AreaId                 int     `json:"area_id"`
	BusinessCategoryId     int     `json:"business_category_id"`
	CompanyPhoneNumber     string  `json:"company_phone_number"`
	IncentiveGroupId       int     `json:"incentive_group_id"`
	CompanyNoOfStall       float64 `json:"company_no_of_stall"`
	CompanyCode            string  `json:"company_code"`
	CompanyTypeId          int     `json:"company_type_id"`
	CompanyFaxNumber       string  `json:"company_fax_number"`
	AftersalesAreaId       int     `json:"aftersales_area_id"`
	CompanyDealerKiaCode   string  `json:"company_dealer_kia_code"`
	CompanyTypeSellingId   int     `json:"company_type_selling_id"`
	CompanyAbbreviation    string  `json:"company_abbreviation"`
	CompanyEmail           string  `json:"company_email"`
	VatCompanyId           int     `json:"vat_company_id"`
	CompanyOwnershipId     int     `json:"company_ownership_id"`
	CompanyOfficeAddressId int     `json:"company_office_address_id"`
}

type CompanyParams struct {
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
	CompanyId   int    `json:"company_id"`
	CompanyCode string `json:"company_code"`
	CompanyName string `json:"company_name"`
	SortBy      string `json:"sort_by"`
	SortOf      string `json:"sort_of"`
}

func GetAllCompany(params CompanyParams) ([]CompanyMasterDetailResponse, *exceptions.BaseErrorResponse) {
	var getCompany []CompanyMasterDetailResponse
	if params.Limit == 0 {
		params.Limit = 1000000
	}

	baseURL := config.EnvConfigs.GeneralServiceUrl + "company-list"

	queryParams := fmt.Sprintf("page=%d&limit=%d", params.Page, params.Limit)

	v := reflect.ValueOf(params)
	typeOfParams := v.Type()
	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i).Interface()
		if strVal, ok := value.(string); ok && strVal != "" {
			key := typeOfParams.Field(i).Tag.Get("json")
			value := strings.ReplaceAll(strVal, " ", "%20")
			queryParams += "&" + key + "=" + value
		}
	}

	url := baseURL + "?" + queryParams

	err := utils.CallAPI("GET", url, nil, &getCompany)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve company due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "company service is temporarily unavailable"
		}

		return getCompany, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting company by ID"),
		}
	}

	return getCompany, nil
}

func GetCompanyVat(id int) (VatCompany, *exceptions.BaseErrorResponse) {
	var getCompanyMaster VatCompany
	url := config.EnvConfigs.GeneralServiceUrl + "company-vat-by-id/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &getCompanyMaster)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve company due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "company service is temporarily unavailable"
		}

		return getCompanyMaster, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting company by ID"),
		}
	}
	return getCompanyMaster, nil
}

func GetCompanyDataById(companyId int) (GetCompanyByIdResponses, *exceptions.BaseErrorResponse) {
	var companyResponse GetCompanyByIdResponses
	companyUrl := config.EnvConfigs.GeneralServiceUrl + "company/" + strconv.Itoa(companyId)

	err := utils.CallAPI("GET", companyUrl, nil, &companyResponse)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve company  due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "company service is temporarily unavailable"
		}

		return companyResponse, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting company by ID"),
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
	url := config.EnvConfigs.GeneralServiceUrl + "company-multi-id/" + strIds

	err := utils.CallAPI("GET", url, nil, response)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve company  due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "company service is temporarily unavailable"
		}

		return &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting company by ID"),
		}
	}
	return nil
}

func GetCompanyDataByCode(companyCode string) (GetCompanyByIdResponses, *exceptions.BaseErrorResponse) {
	var companyResponse GetCompanyByIdResponses
	companyUrl := config.EnvConfigs.GeneralServiceUrl + "company/" + companyCode

	err := utils.CallAPI("GET", companyUrl, nil, &companyResponse)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve company  due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "company service is temporarily unavailable"
		}

		return companyResponse, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting company by Code"),
		}
	}

	return companyResponse, nil
}
func IsFTZCompany(companyId int) bool {
	return companyId == 139 //  ID for FTZ company check
}
