package generalserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type CreateAddressRequest struct {
	AddressStreet1 string `json:"address_street_1"`
	AddressStreet2 string `json:"address_street_2"`
	AddressStreet3 string `json:"address_street_3"`
	VillageId      int    `json:"village_id"`
}

type UpdateAddressRequest struct {
	AddressStreet1   string `json:"address_street_1"`
	AddressStreet2   string `json:"address_street_2"`
	AddressStreet3   string `json:"address_street_3"`
	VillageId        int    `json:"village_id"`
	AddressType      string `json:"address_type"`
	AddressLatitude  int    `json:"address_latitude"`
	AddressLongitude int    `json:"address_longitude"`
}

type AddressResponse struct {
	IsActive       bool   `json:"is_active"`
	AddressId      int    `json:"address_id"`
	AddressStreet1 string `json:"address_street_1"`
	AddressStreet2 string `json:"address_street_2"`
	AddressStreet3 string `json:"address_street_3"`
	VillageId      int    `json:"village_id"`
	AddressType    string `json:"address_type"`
}

const (
	errorMsgCreate = "error consuming external API to create address"
	errorMsgUpdate = "error consuming external API to update address"
	errorMsgGet    = "error consuming external API to get address by ID"
	errorMsgMulti  = "error consuming external API to get address by multiple IDs"
)

func CreateAddress(request CreateAddressRequest) (AddressResponse, *exceptions.BaseErrorResponse) {
	var response AddressResponse
	url := config.EnvConfigs.GeneralServiceUrl + "address"
	err := utils.CallAPI("POST", url, request, &response)
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New(errorMsgCreate),
		}
	}
	return response, nil
}

func UpdateAddress(id int, request UpdateAddressRequest) (AddressResponse, *exceptions.BaseErrorResponse) {
	var response AddressResponse
	url := config.EnvConfigs.GeneralServiceUrl + "address/" + strconv.Itoa(id)
	err := utils.CallAPI("PUT", url, request, &response)
	if err != nil {
		return response, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New(errorMsgUpdate),
		}
	}
	return response, nil
}

func GetAddressByID(id int) (AddressResponse, *exceptions.BaseErrorResponse) {
	var response AddressResponse
	url := config.EnvConfigs.GeneralServiceUrl + "address/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &response)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve address due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "address service is temporarily unavailable"
		}

		return response, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting address by ID"),
		}
	}
	return response, nil
}

func GetAddressByMultiId(ids []int, response interface{}) *exceptions.BaseErrorResponse {

	ids = utils.RemoveDuplicateIds(ids)

	idStrings := make([]string, 0, len(ids))
	for _, id := range ids {
		if id != 0 {
			idStrings = append(idStrings, strconv.Itoa(id))
		}
	}
	strIds := "[" + strings.Join(idStrings, ",") + "]"

	url := config.EnvConfigs.GeneralServiceUrl + "address-by-multi-id/" + strIds
	err := utils.CallAPI("GET", url, nil, response)
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New(errorMsgMulti),
		}
	}
	return nil
}
