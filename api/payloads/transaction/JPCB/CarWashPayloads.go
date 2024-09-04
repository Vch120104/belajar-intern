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
	CarWashBayId               *int       `json:"car_wash_bay_id"`
	CarWashBayDescription      *string    `json:"car_wash_bay_description"`
	CarWashStatusId            int        `json:"car_wash_status_id"`
	CarWashStatusDescription   string     `json:"car_wash_status_description"`
	StartTime                  float32    `json:"start_time"`
	EndTime                    float32    `json:"end_time"`
	CarWashPriorityId          int        `json:"car_wash_priority_id"`
	CarWashPriorityDescription string     `json:"car_wash_priority_description"`
}

type CarWashPostResponse struct {
	CarWashId             int       `json:"car_wash_id"`
	CompanyId             int       `json:"company_id"`
	WorkOrderSystemNumber int       `json:"work_order_system_number"`
	BayId                 *int      `json:"car_wash_bay_id"`
	StatusId              int       `json:"car_wash_status_id"`
	PriorityId            int       `json:"car_wash_priority_id"`
	CarWashDate           time.Time `json:"car_wash_date"`
	StartTime             float32   `json:"start_time"`
	EndTime               float32   `json:"end_time"`
	ActualTime            float32   `json:"actual_time"`
}

type CarWashPostRequestProps struct {
	WorkOrderSystemNumber int `json:"work_order_system_number" validate:"required"`
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

type CarWashUpdatePriorityRequest struct {
	WorkOrderSystemNumber int `json:"work_order_system_number"`
	CarWashStatusId       int `json:"car_wash_status_id"`
	CarWashPriorityId     int `json:"car_wash_priority_id"`
}

type CarWashPriorityDropDownResponse struct {
	CarWashPriorityId          int    `json:"car_wash_priority_id"`
	CarWashPriorityCode        string `json:"car_wash_priority_code"`
	CarWashPriorityDescription string `json:"car_wash_priority_description"`
	IsActive                   bool   `json:"is_active"`
}

type CarWashErrorDetail struct {
	WorkOrderSystemNumber    int    `json:"work_order_system_number"`
	WorkOrderDocumentNumber  string `json:"work_order_document_number"`
	CarWashBayDescription    string `json:"car_wash_bay_description"`
	CarWashStatusId          int    `json:"car_wash_status_id"`
	CarWashStatusDescription string `json:"car_wash_status_description"`
}

type CarWashWorkOrder struct {
	CarWash           bool `json:"car_wash"`
	CompanyId         int  `json:"company_id"`
	WorkOrderStatusId int  `json:"work_order_status_id"`
}

type CarWashCompanyResponse struct {
	CompanyId        int              `json:"company_id"`
	CompanyCode      string           `json:"company_code"`
	CompanyName      string           `json:"company_name"`
	CompanyType      string           `json:"company_type"`
	CompanyReference CompanyReference `json:"company_reference"`
}

type CompanyReference struct {
	UseJPCB *bool `json:"use_jpcb"`
}
