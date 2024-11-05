package salesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

type VehicleResponse struct {
	VehicleId           int             `json:"vehicle_id"`
	VehicleCode         string          `json:"vehicle_chassis_number"`
	VehicleEngineNumber string          `json:"vehicle_engine_number"`
	VehicleTnkb         string          `json:"vehicle_registration_certificate_tnkb"`
	VehicleProduction   json.RawMessage `json:"vehicle_production_year"`
	VehicleLastKm       json.RawMessage `json:"vehicle_last_km"`
	VehicleBrandId      int             `json:"vehicle_brand_id"`
	VehicleModelId      int             `json:"vehicle_model_id"`
	VehicleModelVariant string          `json:"model_variant_colour_description"`
	VehicleVariantId    int             `json:"vehicle_variant_id"`
	VehicleColourId     int             `json:"vehicle_colour_id"`
	VehicleOwner        string          `json:"vehicle_registration_certificate_owner_name"`
}

type VehicleArrayResponse []struct {
	VehicleId           int             `json:"vehicle_id"`
	VehicleCode         string          `json:"vehicle_chassis_number"`
	VehicleEngineNumber string          `json:"vehicle_engine_number"`
	VehicleTnkb         string          `json:"vehicle_registration_certificate_tnkb"`
	VehicleProduction   json.RawMessage `json:"vehicle_production_year"`
	VehicleLastKm       json.RawMessage `json:"vehicle_last_km"`
	VehicleBrandId      int             `json:"vehicle_brand_id"`
	VehicleModelId      int             `json:"vehicle_model_id"`
	VehicleModelVariant string          `json:"model_variant_colour_description"`
	VehicleVariantId    int             `json:"vehicle_variant_id"`
	VehicleColourId     int             `json:"vehicle_colour_id"`
	VehicleOwner        string          `json:"vehicle_registration_certificate_owner_name"`
}

type VehicleUpdate struct {
	VehicleLastKm          int       `json:"vehicle_last_km"`
	VehicleLastServiceDate time.Time `json:"vehicle_last_service_date"`
}

func GetVehicleByChassisNumber(chassis string) (VehicleResponse, *exceptions.BaseErrorResponse) {
	var vehicle VehicleResponse
	url := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + chassis
	errVehicle := utils.CallAPI("GET", url, nil, &vehicle)
	if errVehicle != nil {
		return vehicle, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error consume external vehicle api",
			Err:        errors.New("error consume external vehicle api"),
		}
	}
	return vehicle, nil
}

func UpdateVehicle(id int, request VehicleUpdate) *exceptions.BaseErrorResponse {
	url := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(id)
	errVehicle := utils.CallAPI("PUT", url, request, nil)
	if errVehicle != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error consume external vehicle api",
			Err:        errors.New("error consume external vehicle api"),
		}
	}
	return nil
}

func GetVehicleById(id int) (VehicleArrayResponse, *exceptions.BaseErrorResponse) {
	var vehicle VehicleArrayResponse

	baseURL := config.EnvConfigs.SalesServiceUrl + "vehicle-master"

	params := struct {
		Page      int `json:"page"`
		Limit     int `json:"limit"`
		VehicleID int `json:"vehicle_id"`
	}{
		Page:      0,
		Limit:     1000000000,
		VehicleID: id,
	}

	err := utils.GetArray(baseURL, params, &vehicle)
	if err != nil {
		return vehicle, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error consuming external vehicle API",
			Err:        errors.New("error consuming external vehicle API: " + err.Error()),
		}
	}

	return vehicle, nil
}
