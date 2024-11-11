package salesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type VehicleParams struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	VehicleID int `json:"vehicle_id"`
}

type VehicleResponse struct {
	OrderID                                 int     `json:"order_id"`
	VehicleID                               int     `json:"vehicle_id"`
	VehicleChassisNumber                    string  `json:"vehicle_chassis_number"`
	VehicleRegistrationCertificateTNKB      string  `json:"vehicle_registration_certificate_tnkb"`
	VehicleServiceBookingNumber             string  `json:"vehicle_service_booking_number"`
	VehicleRegistrationCertificateOwnerName string  `json:"vehicle_registration_certificate_owner_name"`
	ModelVariantColourDescription           string  `json:"model_variant_colour_description"`
	VehicleProductionYear                   float64 `json:"vehicle_production_year"`
	VehicleLastServiceDate                  string  `json:"vehicle_last_service_date"`
	VehicleLastKm                           float64 `json:"vehicle_last_km"`
	ColourPoliceName                        string  `json:"colour_police_name"`
	ColourCommercialName                    string  `json:"colour_commercial_name"`
	VehicleBrandID                          int     `json:"vehicle_brand_id"`
	VehicleModelID                          int     `json:"vehicle_model_id"`
	VehicleVariantID                        int     `json:"vehicle_variant_id"`
	VehicleColourID                         int     `json:"vehicle_colour_id"`
	IsActive                                bool    `json:"is_active"`
}

type VehicleListResponse struct {
	TotalRows  int               `json:"total_rows"`
	TotalPages int               `json:"total_pages"`
	Data       []VehicleResponse `json:"data"`
}

type VehicleUpdate struct {
	VehicleLastKm          int       `json:"vehicle_last_km"`
	VehicleLastServiceDate time.Time `json:"vehicle_last_service_date"`
}

// Functions

func GetVehicleByChassisNumber(chassis string) (VehicleResponse, *exceptions.BaseErrorResponse) {
	var vehicleResponse VehicleResponse
	url := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + chassis
	errVehicle := utils.CallAPI("GET", url, nil, &vehicleResponse)
	if errVehicle != nil {
		return VehicleResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error consuming external vehicle API",
			Err:        errors.New("error consuming external vehicle API"),
		}
	}

	return vehicleResponse, nil
}

func UpdateVehicle(id int, request VehicleUpdate) *exceptions.BaseErrorResponse {
	url := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(id)
	errVehicle := utils.CallAPI("PUT", url, request, nil)
	if errVehicle != nil {
		return &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error consuming external vehicle API",
			Err:        errors.New("error consuming external vehicle API"),
		}
	}
	return nil
}

func GetVehicleById(id int) (VehicleResponse, *exceptions.BaseErrorResponse) {
	var vehicleResponse VehicleListResponse

	baseURL := config.EnvConfigs.SalesServiceUrl + "vehicle-master"

	params := VehicleParams{
		Page:      0,
		Limit:     100000,
		VehicleID: id,
	}

	finalURL := fmt.Sprintf("%s?page=%d&limit=%d&vehicle_id=%d", baseURL, params.Page, params.Limit, params.VehicleID)
	//log.Printf("Final URL: %s", finalURL)

	// Make the GET request
	err := utils.GetArray(finalURL, nil, &vehicleResponse) // Ensure you pass nil if no body is needed
	if err != nil {
		return VehicleResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "error consuming external vehicle API: " + err.Error(),
		}
	}

	if len(vehicleResponse.Data) == 0 {
		return VehicleResponse{}, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusNotFound,
			Message:    "vehicle not found",
		}
	}

	return vehicleResponse.Data[0], nil
}
