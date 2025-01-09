package salesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const customTimeLayout = "2006-01-02T15:04:05.999999"

type VehicleParams struct {
	Page                               int    `json:"page"`
	Limit                              int    `json:"limit"`
	VehicleID                          int    `json:"vehicle_id"`
	VehicleChassisNumber               string `json:"vehicle_chassis_number"`
	VehicleRegistrationCertificateTNKB string `json:"vehicle_registration_certificate_tnkb"`
	SortBy                             string `json:"sort_by"`
	SortOf                             string `json:"sort_of"`
}

type Address struct {
	Street1     string `json:"stnk_address_street_1"`
	Street2     string `json:"stnk_address_street_2"`
	Street3     string `json:"stnk_address_street_3"`
	VillageCode string `json:"stnk_village_code"`
	VillageDesc string `json:"stnk_village_desc"`
	Subdistrict string `json:"stnk_subdistrict"`
	Province    string `json:"stnk_province"`
	City        string `json:"stnk_city"`
	ZipCode     string `json:"stnk_zip_code"`
}

type VehicleMaster struct {
	VehicleID                   int       `json:"vehicle_id"`
	IsActive                    bool      `json:"is_active"`
	VehicleChassisNumber        string    `json:"vehicle_chassis_number"`
	VehicleEngineNumber         string    `json:"vehicle_engine_number"`
	VehicleBrandID              int       `json:"vehicle_brand_id"`
	VehicleModelID              int       `json:"vehicle_model_id"`
	VehicleVariantID            int       `json:"vehicle_variant_id"`
	VehicleColourID             int       `json:"vehicle_colour_id"`
	OptionID                    int       `json:"option_id"`
	VehicleTransmissionID       int       `json:"vehicle_transmission_id"`
	VehicleWheelDriveID         int       `json:"vehicle_wheel_drive_id"`
	VehicleCylinderID           int       `json:"vehicle_cylinder_id"`
	VehicleFuelID               int       `json:"vehicle_fuel_id"`
	VehicleProductionYear       int       `json:"vehicle_production_year"`
	VehicleAssemblyYear         int       `json:"vehicle_assembly_year"`
	VehicleVIN2                 string    `json:"vehicle_vin2"`
	VehicleKeyCode              string    `json:"vehicle_key_code"`
	VehicleRadiotapeCode        string    `json:"vehicle_radiotape_code"`
	VehicleServiceBookingNumber string    `json:"vehicle_service_booking_number"`
	VehicleHandoverDocumentDate time.Time `json:"vehicle_handover_document_date"`
	VehicleDealerStock          string    `json:"vehicle_dealer_stock"`
	VehicleLastKM               int       `json:"vehicle_last_km"`
	VehicleStatusID             int       `json:"vehicle_status_id"`
	VehicleIsGreyMarket         bool      `json:"vehicle_is_grey_market"`
	AccessoriesOptionCode       string    `json:"accessories_option_code"`
	AccessoriesOptionName       string    `json:"accessories_option_name"`
	VehicleVariantDescription   string    `json:"vehicle_variant_description"`
	VehicleVariantCode          string    `json:"vehicle_variant_code"`
	VehicleLastServiceDate      time.Time `json:"vehicle_last_service_date"`
	VehicleColourCode           string    `json:"vehicle_colour_code"`
	VehicleColourCommercialName string    `json:"vehicle_colour_commercial_name"`
}

type VehicleSTNK struct {
	VehicleRegistrationCertificateTNKB       string    `json:"vehicle_registration_certificate_tnkb"`
	VehicleRegistrationCertificateNumber     string    `json:"vehicle_registration_certificate_number"`
	VehicleRegistrationCertificateValidDate  time.Time `json:"vehicle_registration_certificate_valid_date"`
	VehicleRegistrationCertificateCode       string    `json:"vehicle_registration_certificate_code"`
	VehicleRegistrationCertificateOwnerName  string    `json:"vehicle_registration_certificate_owner_name"`
	DetailAddress                            Address   `json:"detail_address"`
	StnkVillageID                            int       `json:"stnk_village_id"`
	VehicleRegistrationCertificateTNKBColour string    `json:"vehicle_registration_certificate_tnkb_colour"`
	VehicleBPKBNumber                        string    `json:"vehicle_bpkb_number"`
	VehicleBPKBDate                          time.Time `json:"vehicle_bpkb_date"`
	PoliceInvoiceNumber                      string    `json:"police_invoice_number"`
}

type VehicleCurrentUser struct {
	UserCustomerID    int     `json:"user_customer_id"`
	UserCustomerCode  string  `json:"user_customer_code"`
	UserCustomerName  string  `json:"user_customer_name"`
	DetailUserAddress Address `json:"detail_user_address"`
}

type VehicleBilling struct {
	BillingCustomerID    int     `json:"billing_customer_id"`
	BillingCustomerCode  string  `json:"billing_customer_code"`
	BillingCustomerName  string  `json:"billing_customer_name"`
	DetailBillingAddress Address `json:"detail_billing_address"`
}

type VehicleContractService struct {
	ContractID            int       `json:"contract_id"`
	ContractServiceNumber string    `json:"contract_service_number"`
	ContractServiceDate   time.Time `json:"contract_service_date"`
	ContractDealer        int       `json:"contract_dealer"`
}

type VehicleInsurance struct {
	InsuranceCompany      string `json:"insurance_company"`
	InsurancePolicyNumber string `json:"insurance_policy_number"`
	InsuranceEndDate      string `json:"insurance_end_date"`
	InsuranceEraNumber    string `json:"insurance_era_number"`
}

type VehicleReceiving struct {
	VehicleWRSNumber           string    `json:"vehicle_wrs_number"`
	VehicleWRSDate             time.Time `json:"vehicle_wrs_date"`
	VehicleBPUNumber           string    `json:"vehicle_bpu_number"`
	VehicleBPUDate             time.Time `json:"vehicle_bpu_date"`
	VehicleDeliveryOrderNumber string    `json:"vehicle_delivery_order_number"`
	VehicleDeliveryOrderDate   time.Time `json:"vehicle_delivery_order_date"`
	VehicleSJNumber            string    `json:"vehicle_sj_number"`
	VehicleSJDate              time.Time `json:"vehicle_sj_date"`
	VehicleExpedition          string    `json:"vehicle_expedition"`
}

type VehicleLastStatus struct {
	BookedBy                        string    `json:"booked_by"`
	BookingNumber                   string    `json:"booking_number"`
	BookingDate                     time.Time `json:"booking_date"`
	VehicleSalesOrderDocumentNumber string    `json:"vehicle_sales_order_document_number"`
	VehicleSalesOrderDate           time.Time `json:"vehicle_sales_order_date"`
	VehicleLastFakturPolisiNumber   string    `json:"vehicle_last_faktur_polisi_number"`
	VehicleLastSTNKNumber           string    `json:"vehicle_last_stnk_number"`
	VehicleLastSTNKDate             time.Time `json:"vehicle_last_stnk_date"`
	VehicleLastBPKBNumber           string    `json:"vehicle_last_bpkb_number"`
	VehicleLastBPKBDate             time.Time `json:"vehicle_last_bpkb_date"`
	VehicleLastNikNumber            string    `json:"vehicle_last_nik_number"`
	BPKDocumentNumber               string    `json:"bpk_document_number"`
	BPKDate                         time.Time `json:"bpk_date"`
	VehicleUnitDelivery             string    `json:"vehicle_unit_delivery"`
	VehicleUnitStockStatus          string    `json:"vehicle_unit_stock_status"`
	WarehouseID                     int       `json:"warehouse_id"`
	WarehouseName                   string    `json:"warehouse_name"`
	VehicleUnitLotNumber            string    `json:"vehicle_unit_lot_number"`
}

type VehiclePurchase struct {
	VehiclePurchaseDate    time.Time `json:"vehicle_purchase_date"`
	VehiclePurchaseOrderID int       `json:"vehicle_purchase_order_id"`
	VehiclePurchaseDealer  int       `json:"vehicle_purchase_dealer"`
	VehicleCustomerCode    string    `json:"vehicle_customer_code"`
	VehiclePurchasePrice   int       `json:"vehicle_purchase_price"`
}

type VehicleResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       struct {
		Master          VehicleMaster          `json:"master"`
		STNK            VehicleSTNK            `json:"stnk"`
		CurrentUser     VehicleCurrentUser     `json:"current_user"`
		Billing         VehicleBilling         `json:"billing"`
		ContractService VehicleContractService `json:"contract_service"`
		Insurance       VehicleInsurance       `json:"insurance"`
		Receiving       VehicleReceiving       `json:"receiving"`
		LastStatus      VehicleLastStatus      `json:"last_status"`
		Purchase        VehiclePurchase        `json:"purchase"`
	} `json:"data"`
}

type VehicleListResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	PageLimit  int    `json:"page_limit"`
	Page       int    `json:"page"`
	TotalRows  int    `json:"total_rows"`
	TotalPages int    `json:"total_pages"`
	Data       []struct {
		VehicleID                   int    `json:"vehicle_id"`
		IsActive                    bool   `json:"is_active"`
		VehicleChassisNumber        string `json:"vehicle_chassis_number"`
		VehicleBrandID              int    `json:"vehicle_brand_id"`
		VehicleModelID              int    `json:"vehicle_model_id"`
		VehicleVariantID            int    `json:"vehicle_variant_id"`
		VehicleColourID             int    `json:"vehicle_colour_id"`
		VehicleProductionYear       int    `json:"vehicle_production_year"`
		VehicleLastKm               int    `json:"vehicle_last_km"`
		VehicleLastServiceDate      string `json:"vehicle_last_service_date"`
		VehicleColourCommercialName string `json:"vehicle_colour_commercial_name"`
	} `json:"data"`
}

type VehicleUpdate struct {
	VehicleLastKm          int       `json:"vehicle_last_km"`
	VehicleLastServiceDate time.Time `json:"vehicle_last_service_date"`
}

func createVehicleURL(baseURL string, params VehicleParams) string {
	queryParams := fmt.Sprintf("page=%d&limit=%d", params.Page, params.Limit)

	v := reflect.ValueOf(params)
	typeOfParams := v.Type()
	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i).Interface()
		if strVal, ok := value.(string); ok && strVal != "" {
			key := typeOfParams.Field(i).Tag.Get("json")
			queryParams += "&" + key + "=" + strings.ReplaceAll(strVal, " ", "%20")
		}
	}

	return baseURL + "?" + queryParams
}

func GetVehicleByChassisNumber(chassis string) (VehicleResponse, *exceptions.BaseErrorResponse) {
	var vehicleResponse VehicleResponse
	url := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + chassis
	errVehicle := utils.CallAPI("GET", url, nil, &vehicleResponse)
	if errVehicle != nil {
		status := http.StatusBadGateway
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

func GetAllVehicle(params VehicleParams) ([]VehicleListResponse, *exceptions.BaseErrorResponse) {
	var getVehicle []VehicleListResponse
	if params.Limit == 0 {
		params.Limit = 1000000
	}

	baseURL := config.EnvConfigs.GeneralServiceUrl + "vehicle-master"
	url := createVehicleURL(baseURL, params)

	err := utils.CallAPI("GET", url, nil, &getVehicle)
	if err != nil {
		status := http.StatusBadGateway
		message := "Failed to retrieve vehicles due to an external service error"
		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "unit vehicle is temporarily unavailable"
		}
		return getVehicle, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting vehicles"),
		}
	}

	return getVehicle, nil
}

func parseCustomTime(timeStr string) (time.Time, error) {
	t, err := time.Parse(customTimeLayout, timeStr)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func GetVehicleById(id int) (VehicleResponse, *exceptions.BaseErrorResponse) {
	var vehicleData VehicleResponse
	url := config.EnvConfigs.SalesServiceUrl + "vehicle-master/" + strconv.Itoa(id)

	resp, err := http.Get(url)
	if err != nil {
		status := http.StatusBadGateway
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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return vehicleData, &exceptions.BaseErrorResponse{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("Error: Received non-OK HTTP status: %s", resp.Status),
			Err:        fmt.Errorf("received non-OK HTTP status: %s", resp.Status),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return vehicleData, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to read response body",
			Err:        err,
		}
	}

	err = json.Unmarshal(body, &vehicleData)
	if err != nil {
		return vehicleData, &exceptions.BaseErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to decode vehicle data from external service response",
			Err:        err,
		}
	}

	log.Printf("Full Vehicle Data: %+v", vehicleData)

	fieldsToParse := []struct {
		name    string
		timeStr string
	}{
		// Master Fields
		{"vehicle_handover_document_date", vehicleData.Data.Master.VehicleHandoverDocumentDate.Format(customTimeLayout)},
		{"vehicle_last_service_date", vehicleData.Data.Master.VehicleLastServiceDate.Format(customTimeLayout)},

		// STNK Fields
		{"vehicle_registration_certificate_valid_date", vehicleData.Data.STNK.VehicleRegistrationCertificateValidDate.Format(customTimeLayout)},
		{"vehicle_bpkb_date", vehicleData.Data.STNK.VehicleBPKBDate.Format(customTimeLayout)},

		// Contract Service Fields
		{"contract_service_date", vehicleData.Data.ContractService.ContractServiceDate.Format(customTimeLayout)},

		// Receiving Fields
		{"vehicle_wrs_date", vehicleData.Data.Receiving.VehicleWRSDate.Format(customTimeLayout)},
		{"vehicle_sj_date", vehicleData.Data.Receiving.VehicleSJDate.Format(customTimeLayout)},

		// Last Status Fields
		{"vehicle_sales_order_date", vehicleData.Data.LastStatus.VehicleSalesOrderDate.Format(customTimeLayout)},
		{"vehicle_last_stnk_date", vehicleData.Data.LastStatus.VehicleLastSTNKDate.Format(customTimeLayout)},
		{"vehicle_last_bpkb_date", vehicleData.Data.LastStatus.VehicleLastBPKBDate.Format(customTimeLayout)},
		{"bpk_date", vehicleData.Data.LastStatus.BPKDate.Format(customTimeLayout)},

		// Purchase Fields
		{"vehicle_purchase_date", vehicleData.Data.Purchase.VehiclePurchaseDate.Format(customTimeLayout)},

		// Insurance Fields
		{"insurance_end_date", vehicleData.Data.Insurance.InsuranceEndDate},
	}

	for _, field := range fieldsToParse {
		parsedTime, err := parseCustomTime(field.timeStr)
		if err != nil {
			log.Printf("Error parsing %s: %v", field.name, err)
		} else {
			log.Printf("Parsed %s: %v", field.name, parsedTime)
		}
	}

	if vehicleData.Data.Master.VehicleID == 0 {
		log.Printf("Warning: VehicleID is empty for Vehicle ID: %d", vehicleData.Data.Master.VehicleID)
	} else {
		log.Printf("Vehicle ID: %d", vehicleData.Data.Master.VehicleID)
	}

	return vehicleData, nil
}
