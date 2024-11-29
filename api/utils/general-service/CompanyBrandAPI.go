package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type CompanyBrandParams struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type CompanyBrandResponse struct {
	IsActive         bool   `json:"is_active"`
	CompanyBrandId   int    `json:"company_brand_id"`
	BrandId          int    `json:"brand_id"`
	BrandCode        string `json:"brand_code"`
	BrandName        string `json:"brand_name"`
	BusinessTypeName string `json:"business_type_name"`
	GenerateAccPo    bool   `json:"generate_acc_po"`
}

type CompanyBrandByCompanyResponse struct {
	StatusCode int                    `json:"status_code"`
	Message    string                 `json:"message"`
	Page       int                    `json:"page"`
	PageLimit  int                    `json:"limit"`
	TotalRows  int                    `json:"total_rows"`
	TotalPages int                    `json:"total_pages"`
	Data       []CompanyBrandResponse `json:"data"`
}

func GetCompanyBrandByCompanyPagination(id int, params CompanyBrandParams) (CompanyBrandByCompanyResponse, *exceptions.BaseErrorResponse) {
	var getCompanyBrand CompanyBrandByCompanyResponse
	if params.Limit == 0 {
		params.Limit = 100000
	}
	url := config.EnvConfigs.GeneralServiceUrl + "company-brand-list/" + strconv.Itoa(id) + "?page=" + strconv.Itoa(params.Page) + "&limit=" + strconv.Itoa(params.Limit)

	err := utils.GetArray(url, nil, &getCompanyBrand)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve company brand due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "company brand service is temporarily unavailable"
		}

		return getCompanyBrand, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting company brand"),
		}
	}
	return getCompanyBrand, nil
}
