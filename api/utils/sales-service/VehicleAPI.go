package salesserviceapiutils

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
	"time"
)

type VehicleParams struct {
	Page                               int    `json:"page"`
	Limit                              int    `json:"limit"`
	VehicleID                          int    `json:"vehicle_id"`
	VehicleChassisNumber               string `json:"vehicle_chassis_number"`
	VehicleRegistrationCertificateTNKB string `json:"vehicle_registration_certificate_tnkb"`
	SortBy                             string `json:"sort_by"`
	SortOf                             string `json:"sort_of"`
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
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve vehicle due to an external service error"

		if errors.Is(errVehicle, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "vehicle service is temporarily unavailable"
		}

		return vehicleResponse, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting vehicle by ID"),
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

func GetAllVehicle(params VehicleParams) ([]VehicleResponse, *exceptions.BaseErrorResponse) {
	var getVehicle []VehicleResponse
	if params.Limit == 0 {
		params.Limit = 1000000
	}

	baseURL := config.EnvConfigs.GeneralServiceUrl + "vehicle-master"

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

	err := utils.CallAPI("GET", url, nil, &getVehicle)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve vehicle due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "unit vehicle is temporarily unavailable"
		}

		return getVehicle, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting vehicle by ID"),
		}
	}

	return getVehicle, nil
}

func GetVehicleById(id int) (VehicleResponse, *exceptions.BaseErrorResponse) {
	var vehicleData VehicleResponse
	url := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(id)
	err := utils.CallAPI("GET", url, nil, &vehicleData)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve vehicle due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "brand service is temporarily unavailable"
		}

		return vehicleData, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting vehicle by ID"),
		}
	}
	return vehicleData, nil
}
