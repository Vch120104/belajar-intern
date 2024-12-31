package salesserviceapiutils

import (
	"after-sales/api/config"
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type VehicleSalesOrderResponse struct {
	VehicleSalesOrderStatus             int    `json:"vehicle_sales_order_status"`
	SalesRepresentativeId               int    `json:"sales_representative_id"`
	SalesPartnersId                     int    `json:"sales_partners_id"`
	VehicleSalesOrderSystemNumber       int    `json:"vehicle_sales_order_system_number"`
	VehicleSalesOrderDocumentNumber     string `json:"vehicle_sales_order_document_number"`
	VehicleSalesOrderStageDate          string `json:"vehicle_sales_order_stage_date"`
	ProspectSystemNumber                int    `json:"prospect_system_number"`
	DealerRepresentative                int    `json:"dealer_representative"`
	Keyword                             string `json:"keyword"`
	BuyingTypeId                        int    `json:"buying_type_id"`
	OrderByCustomerId                   int    `json:"order_by_customer_id"`
	CustomerCode                        string `json:"customer_code"`
	CustomerIdType                      int    `json:"customer_id_type"`
	CustomerIdNumber                    string `json:"customer_id_number"`
	CustomerTitlePrefix                 string `json:"customer_title_prefix"`
	CustomerName                        string `json:"customer_name"`
	CustomerTitleSuffix                 string `json:"customer_title_suffix"`
	CustomerGender                      int    `json:"customer_gender"`
	CustomerAddressStreet1              string `json:"customer_address_street_1"`
	CustomerAddressStreet2              string `json:"customer_address_street_2"`
	CustomerAddressStreet3              string `json:"customer_address_street_3"`
	CustomerVillageId                   int    `json:"customer_village_id"`
	CustomerPhoneNo                     string `json:"customer_phone_no"`
	CustomerHomeFaxNo                   string `json:"customer_home_fax_no"`
	CustomerMobileNo                    string `json:"customer_mobile_no"`
	CustomerEmailAddress                string `json:"customer_email_address"`
	ClientTypeId                        int    `json:"client_type_id"`
	CustomerCategoryId                  int    `json:"customer_category_id"`
	UnitTransactionTypeId               int    `json:"unit_transaction_type_id"`
	BusinessTypeId                      int    `json:"business_type_id"`
	BusinessGroupId                     int    `json:"business_group_id"`
	BusinessWebsite                     string `json:"business_website"`
	CorporationPurchaseOrderNo          string `json:"corporation_purchase_order_no"`
	CorporationPurchaseOrderDate        string `json:"corporation_purchase_order_date"`
	ContactName                         string `json:"contact_name"`
	ContactGenderId                     int    `json:"contact_gender_id"`
	ContactPhoneNumber                  string `json:"contact_phone_number"`
	ContactEmailAddress                 string `json:"contact_email_address"`
	TaxInvoiceTypeId                    int    `json:"tax_invoice_type_id"`
	VatRegistrationNumber               string `json:"vat_registration_number"`
	VatRegistrationDate                 string `json:"vat_registration_date"`
	VatName                             string `json:"vat_name"`
	VatAddressStreet1                   string `json:"vat_address_street_1"`
	VatAddressStreet2                   string `json:"vat_address_street_2"`
	VatAddressStreet3                   string `json:"vat_address_street_3"`
	VatVillageId                        int    `json:"vat_village_id"`
	VatPkpNumber                        string `json:"vat_pkp_number"`
	VatPkpDate                          string `json:"vat_pkp_date"`
	StnkCustomerId                      int    `json:"stnk_customer_id"`
	StnkCustomerCode                    string `json:"stnk_customer_code"`
	StnkCustomerIdType                  int    `json:"stnk_customer_id_type"`
	StnkCustomerIdNumber                int    `json:"stnk_customer_id_number"`
	StnkCustomerTitlePrefix             string `json:"stnk_customer_title_prefix"`
	StnkCustomerName                    string `json:"stnk_customer_name"`
	StnkCustomerTitleSuffix             string `json:"stnk_customer_title_suffix"`
	StnkCustomerGender                  int    `json:"stnk_customer_gender"`
	StnkCustomerAddressStreet1          string `json:"stnk_customer_address_street_1"`
	StnkCustomerAddressStreet2          string `json:"stnk_customer_address_street_2"`
	StnkCustomerAddressStreet3          string `json:"stnk_customer_address_street_3"`
	StnkCustomerVillageId               int    `json:"stnk_customer_village_id"`
	StnkCustomerPhoneNo                 string `json:"stnk_customer_phone_no"`
	StnkCustomerHomeFaxNo               string `json:"stnk_customer_home_fax_no"`
	StnkCustomerEmailAddress            string `json:"stnk_customer_email_address"`
	CorrespondentCustomerId             int    `json:"correspondent_customer_id"`
	CorrespondentCustomerCode           string `json:"correspondent_customer_code"`
	CorrespondentCustomerIdType         int    `json:"correspondent_customer_id_type"`
	CorrespondentCustomerIdNumber       string `json:"correspondent_customer_id_number"`
	CorrespondentCustomerTitlePrefix    string `json:"correspondent_customer_title_prefix"`
	CorrespondentCustomerName           string `json:"correspondent_customer_name"`
	CorrespondentCustomerTitleSuffix    string `json:"correspondent_customer_title_suffix"`
	CorrespondentCustomerGender         int    `json:"correspondent_customer_gender"`
	CorrespondentCustomerAddressStreet1 string `json:"correspondent_customer_address_street_1"`
	CorrespondentCustomerAddressStreet2 string `json:"correspondent_customer_address_street_2"`
	CorrespondentCustomerAddressStreet3 string `json:"correspondent_customer_address_street_3"`
	CorrespondentCustomerVillageId      int    `json:"correspondent_customer_village_id"`
	CorrespondentCustomerPhoneNo        string `json:"correspondent_customer_phone_no"`
	CorrespondentCustomerHomeFaxNo      string `json:"correspondent_customer_home_fax_no"`
	CorrespondentCustomerEmailAddress   string `json:"correspondent_customer_email_address"`
	ProspectSourceId                    int    `json:"prospect_source_id"`
	ProspectDate                        string `json:"prospect_date"`
	Note                                string `json:"note"`
	PriceCode                           int    `json:"price_code"`
	IncentiveStatus                     bool   `json:"incentive_status"`
	PlateColorId                        int    `json:"plate_color_id"`
	BodyType                            int    `json:"body_type"`
}

func GetVehicleSalesOrderById(id int) (VehicleSalesOrderResponse, *exceptions.BaseErrorResponse) {
	// https://testing-backendims.indomobil.co.id/sales-service/v1/vehicle-sales-order/53
	Url := config.EnvConfigs.SalesServiceUrl + "vehicle-sales-order/" + strconv.Itoa(id)
	response := []VehicleSalesOrderResponse{}
	err := utils.CallAPI("GET", Url, nil, &response)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve unit variant due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "unit variant service is temporarily unavailable"
		}

		return response[0], &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting unit variant by brand"),
		}
	}
	return response[0], nil
}
