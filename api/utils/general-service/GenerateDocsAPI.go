package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type BrandDocResponse struct {
	BrandId           int    `json:"brand_id"`
	BrandCode         string `json:"brand_code"`
	BrandName         string `json:"brand_name"`
	BrandAbbreviation string `json:"brand_abbreveation"`
}

func GetBrandGenerateDoc(id int) (BrandDocResponse, *exceptions.BaseErrorResponse) {
	var brandDoc BrandDocResponse
	url := config.EnvConfigs.SalesServiceUrl + "unit-brand/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, &brandDoc, nil)
	if err != nil {
		return brandDoc, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching brand doc by code",
			Err:        errors.New("error consuming external API while fetching brand doc by code"),
		}
	}
	return brandDoc, nil
}
