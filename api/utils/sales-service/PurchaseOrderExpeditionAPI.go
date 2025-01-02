package salesserviceapiutils

import (
	"after-sales/api/exceptions"
	"after-sales/api/utils"
	"errors"
	"net/http"
	"strconv"
)

type PurchaseOrderExpeditionResponse struct {
	CompanyId                                 int         `json:"company_id"`
	WarehouseContactPerson                    string      `json:"warehouse_contact_person"`
	UnitPurchaseOrderExpeditionStatusId       int         `json:"unit_purchase_order_expedition_status_id"`
	SupplierId                                int         `json:"supplier_id"`
	BillCodeId                                int         `json:"bill_code_id"`
	ContactPersonMobilePhone                  string      `json:"contact_person_mobile_phone"`
	BrandId                                   int         `json:"brand_id"`
	PickupWarehouseId                         int         `json:"pickup_warehouse_id"`
	FareTypeId                                int         `json:"fare_type_id"`
	EstimatedPickupDate                       string      `json:"estimated_pickup_date"`
	DeliveryWarehouseId                       int         `json:"delivery_warehouse_id"`
	EstimatedPickupTime                       string      `json:"estimated_pickup_time"`
	UnitPurchaseOrderExpeditionDate           string      `json:"unit_purchase_order_expedition_date"`
	PickupAddressId                           int         `json:"pickup_address_id"`
	TermOfPaymentId                           int         `json:"term_of_payment_id"`
	UnitPurchaseOrderExpeditionDocumentNumber string      `json:"unit_purchase_order_expedition_document_number"`
	PurchaseOrderTypeExpeditionId             int         `json:"purchase_order_type_expedition_id"`
	DeliveryAddressId                         int         `json:"delivery_address_id"`
	UnitPurchaseOrderExpeditionRemark         string      `json:"unit_purchase_order_expedition_remark"`
	UnitPurchaseOrderExpeditionSystemNumber   int         `json:"unit_purchase_order_expedition_system_number"`
	BpkSystemNumber                           interface{} `json:"bpk_system_number"`
}

const PurchaseOrderExpeditionBaseUrl = "unit-purchase-order-expedition/"

func GetPurchaseOrderExpeditionById(id int) (PurchaseOrderExpeditionResponse, *exceptions.BaseErrorResponse) {
	response := PurchaseOrderExpeditionResponse{}
	PurchaseOrderExpeditionUrl := PurchaseOrderExpeditionBaseUrl + strconv.Itoa(id)
	err := utils.CallAPI("GET", PurchaseOrderExpeditionUrl, nil, &response)
	if err != nil {
		status := http.StatusBadGateway // Default to 502
		message := "Failed to retrieve purchase order expedition due to an external service error"

		if errors.Is(err, utils.ErrServiceUnavailable) {
			status = http.StatusServiceUnavailable
			message = "brand service is temporarily unavailable"
		}

		return response, &exceptions.BaseErrorResponse{
			StatusCode: status,
			Message:    message,
			Err:        errors.New("error consuming external API while getting purchase order expedition by code"),
		}
	}
	return response, nil
}
