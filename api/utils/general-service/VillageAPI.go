package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"fmt"
	"net/http"
	"strconv"
)

type VillageByIdResponse struct {
	VillageId      int    `json:"village_id"`
	VillageCode    string `json:"village_code"`
	VillageName    string `json:"village_name"`
	DistrictId     int    `json:"district_id"`
	DistrictName   string `json:"district_name"`
	DistrictCode   string `json:"district_code"`
	CityName       string `json:"city_name"`
	ProvinceName   string `json:"province_name"`
	CountryName    string `json:"country_name"`
	VillageZipCode string `json:"village_zip_code"`
	CityPhoneArea  string `json:"city_phone_area"`
	IsActive       bool   `json:"is_active"`
	CityId         int    `json:"city_id"`
	ProvinceId     int    `json:"province_id"`
}

func GetVillageByID(id int) (VillageByIdResponse, *exceptions.BaseErrorResponse) {
	var getVillage VillageByIdResponse
	url := config.EnvConfigs.GeneralServiceUrl + "village/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &getVillage)
	if err != nil {
		return getVillage, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed to get village by id : %w", err),
		}
	}
	return getVillage, nil
}
