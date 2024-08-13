package transactionjpcbpayloads

import "time"

type CarWashGetAllResponse struct {
	WorkOrderSystemNumber      int        `json:"work_order_system_number"`
	WorkOrderDocumentNumber    string     `json:"work_order_document_number"`
	Model                      string     `json:"model"`
	Color                      string     `json:"color"`
	Tnkb                       string     `json:"tnkb"`
	PromiseTime                *time.Time `json:"promise_time"`
	PromiseDate                *time.Time `json:"promise_date"`
	CarWashBayId               int        `json:"car_wash_bay_id"`
	CarWashStatusId            int        `json:"car_wash_status_id"`
	CarWashStatusDescription   string     `json:"car_wash_status_description"`
	StartTime                  *time.Time `json:"start_time"`
	EndTime                    *time.Time `json:"end_time"`
	CarWashPriorityId          int        `json:"car_wash_priority_id"`
	CarWashPriorityDescription string     `json:"car_wash_priority_description"`
}

type CarWashModelResponse struct {
	ModelId   int    `json:"model_id"`
	ModelCode string `json:"model_code"`
	ModelName string `json:"model_description"`
}

type CarWashVehicleResponse struct {
	VehicleId       int `json:"vehicle_id"`
	VehicleColourId int `json:"vehicle_colour_id"`
}

type CarWashColourResponse struct {
	VariantColourId   int    `json:"colour_id"`
	VariantColourCode string `json:"colour_commercial_name"`
	VariantColourName string `json:"colour_police_name"`
}
