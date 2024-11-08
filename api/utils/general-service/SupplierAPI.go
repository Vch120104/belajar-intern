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

type SupplierMasterVatResponse struct {
	NpwpNumber         string `json:"npwp_no"`
	NpwpDate           string `json:"npwp_date"`
	PkpType            bool   `json:"pkp_type"`
	PkpNumber          string `json:"pkp_no"`
	PkpDate            string `json:"pkp_date"`
	TaxTransactionId   int    `json:"tax_transaction_id"`
	Name               string `json:"name"`
	AddressStreet1     string `json:"address_street_1"`
	AddressStreet2     string `json:"address_street_2"`
	AddressStreet3     string `json:"address_street_3"`
	VillageId          int    `json:"village_id"`
	TaxServiceOfficeId int    `json:"tax_service_office_id"`
}

type SupplierMasterResponse struct {
	SupplierId                int                       `json:"supplier_id"`
	SupplierCode              string                    `json:"supplier_code"`
	SupplierName              string                    `json:"supplier_name"`
	SupplierTypeId            int                       `json:"supplier_type_id"`
	SupplierMasterVatResponse SupplierMasterVatResponse `json:"vat_supplier"`
	ClientTypeId              int                       `json:"client_type_id"`
}

func GetSupplierMasterByCode(code string) (SupplierMasterResponse, *exceptions.BaseErrorResponse) {
	var getSupplierMaster SupplierMasterResponse
	url := config.EnvConfigs.GeneralServiceUrl + "supplier-by-code/" + code

	err := utils.CallAPI("GET", url, nil, &getSupplierMaster)
	if err != nil {
		return getSupplierMaster, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching supplier by code",
			Err:        errors.New("failed to retrieve supplier data from external API by code"),
		}
	}
	return getSupplierMaster, nil
}

func GetSupplierMasterByID(id int) (SupplierMasterResponse, *exceptions.BaseErrorResponse) {
	var getSupplierMaster SupplierMasterResponse
	url := config.EnvConfigs.GeneralServiceUrl + "supplier/" + strconv.Itoa(id)

	err := utils.CallAPI("GET", url, nil, &getSupplierMaster)
	if err != nil {
		return getSupplierMaster, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching supplier by ID",
			Err:        errors.New("failed to retrieve supplier data from external API by ID"),
		}
	}
	return getSupplierMaster, nil
}

func GetSupplierMasterByMultiId(ids []int, abstractType interface{}) *exceptions.BaseErrorResponse {
	ids = utils.RemoveDuplicateIds(ids)
	var nonZeroIds []string

	for _, id := range ids {
		if id != 0 {
			nonZeroIds = append(nonZeroIds, strconv.Itoa(id))
		}
	}

	strIds := "[" + strings.Join(nonZeroIds, ",") + "]"
	url := config.EnvConfigs.GeneralServiceUrl + "supplier-multi-id/" + strIds

	err := utils.CallAPI("GET", url, nil, &abstractType)
	if err != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error fetching suppliers by multiple IDs",
			Err:        errors.New("failed to retrieve supplier data from external API for multiple IDs"),
		}
	}
	return nil
}
