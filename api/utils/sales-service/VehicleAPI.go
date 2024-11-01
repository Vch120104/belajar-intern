package salesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"encoding/json"
	"errors"
	"net/http"
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

func GetVehicleByChassisNumber(chassis string) (VehicleResponse, *exceptions.BaseErrorResponse) {
	var vehicle VehicleResponse
	url := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + chassis
	errVehicle := utils.CallAPI("GET", url, nil, &vehicle)
	if errVehicle != nil {
		return vehicle, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New("error consume external api"),
		}
	}
	return vehicle, nil
}
