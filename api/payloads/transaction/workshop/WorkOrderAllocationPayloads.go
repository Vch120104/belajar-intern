package transactionworkshoppayloads

import "time"

type WorkOrderAllocationAssignRequest struct {
	CompanyId             int       `json:"company_id"`
	TechnicianId          int       `json:"technician_id"`
	ShiftCode             string    `json:"shift_code"`
	StartTime             float64   `json:"start_time"`
	EndTime               float64   `json:"end_time"`
	TotalHour             float64   `json:"total_hour"`
	AvailableSystemNumber int       `json:"available_system_number"`
	ServiceDateTime       time.Time `json:"service_date_time"`
	ForemanId             int       `json:"foreman_id"`
	ReferenceType         string    `json:"reference_type"`
	ReferenceSystemNumber int       `json:"reference_system_number"`
	ReferenceLine         int       `json:"reference_line"`
	Remark                string    `json:"remark"`
}

type WorkOrderAllocationAssignResponse struct {
	CompanyId             int       `json:"company_id"`
	TechnicianId          int       `json:"technician_id"`
	ShiftCode             string    `json:"shift_code"`
	StartTime             float64   `json:"start_time"`
	EndTime               float64   `json:"end_time"`
	TotalHour             float64   `json:"total_hour"`
	AvailableSystemNumber int       `json:"available_system_number"`
	ServiceDateTime       time.Time `json:"service_date_time"`
	ForemanId             int       `json:"foreman_id"`
	ReferenceType         string    `json:"reference_type"`
	ReferenceSystemNumber int       `json:"reference_system_number"`
	ReferenceLine         int       `json:"reference_line"`
	Remark                string    `json:"remark"`
}

type WorkOrderAllocationRequest struct {
	ServiceRequestDate    string  `json:"service_date"`
	BrandId               int     `json:"brand_id"`
	WorkOrderSystemNumber int     `json:"work_order_system_number"`
	ForemanId             int     `json:"foreman_id"`
	ServiceAdvisorId      int     `json:"service_advisor_id"`
	ModelId               int     `json:"model_id"`
	VariantId             int     `json:"variant_id"`
	VehicleId             int     `json:"vehicle_id"`
	CustomerId            int     `json:"customer_id"`
	Frt                   float64 `json:"frt"`
}

type WorkOrderAllocationResponse struct {
	ServiceRequestDate      string  `json:"service_date"`
	BrandId                 int     `json:"brand_id"`
	BrandName               string  `json:"brand_name"`
	WorkOrderSystemNumber   int     `json:"work_order_system_number"`
	WorkOrderDocumentNumber string  `json:"work_order_document_number"`
	ForemanId               int     `json:"foreman_id"`
	ForemanName             string  `json:"foreman_name"`
	ServiceAdvisorId        int     `json:"service_advisor_id"`
	ServiceAdvisorName      string  `json:"service_advisor_name"`
	ModelId                 int     `json:"model_id"`
	ModelName               string  `json:"model_name"`
	VariantId               int     `json:"variant_id"`
	VariantName             string  `json:"variant_name"`
	VehicleId               int     `json:"vehicle_id"`
	VehicleChassisNumber    string  `json:"vehicle_chassis_number"`
	CustomerId              int     `json:"customer_id"`
	CustomerName            string  `json:"customer_name"`
	CustomerBehavior        string  `json:"customer_behavior"`
	Frt                     float64 `json:"frt"`
}

type WorkOrderAllocationDetailRequest struct {
	TechnicianId          int       `json:"technician_id" parent_entity:"trx_work_order_allocation_detail" main_table:"trx_work_order_allocation_detail"`
	WorkOrderSystemNumber int       `json:"work_order_system_number" parent_entity:"trx_work_order_allocation_detail"`
	ShiftCode             string    `json:"shift_code" parent_entity:"trx_work_order_allocation_detail"`
	StartTime             time.Time `json:"start_time" parent_entity:"trx_work_order_allocation_detail"`
	EndTime               time.Time `json:"end_time" parent_entity:"trx_work_order_allocation_detail"`
}

type WorkOrderAllocationDetailResponse struct {
	TechnicianId            int    `json:"technician_id"`
	TechnicianName          string `json:"technician_name"`
	WorkOrderSystemNumber   int    `json:"work_order_system_number"`
	WorkOrderDocumentNumber string `json:"work_order_document_number"`
	ShiftCode               string `json:"shift_code"`
	ServiceStatusId         int    `json:"service_status_id"`
	ServiceStatus           string `json:"service_status"`
	StartTime               string `json:"start_time"`
	EndTime                 string `json:"end_time"`
}

type WorkOrderAllocationGridRequest struct {
	ForemanId   int       `json:"foreman_id"`
	ServiceDate time.Time `json:"service_date"`
	CompanyId   int       `json:"company_id"`
}

type WorkOrderAllocationGridResponse struct {
	ForemanId          int       `json:"foreman_id"`
	ForemanName        string    `json:"foreman_name"`
	CompanyId          int       `json:"company_id"`
	CompanyName        string    `json:"company_name"`
	ServiceDate        time.Time `json:"service_date"`
	TechnicianId       int       `json:"technician_id"`
	TechnicianName     string    `json:"technician_name"`
	ShiftCode          string    `json:"shift_code"`
	TimeAllocation0700 float64   `json:"time_allocation_0700"`
	TimeAllocation0715 float64   `json:"time_allocation_0715"`
	TimeAllocation0730 float64   `json:"time_allocation_0730"`
	TimeAllocation0745 float64   `json:"time_allocation_0745"`
	TimeAllocation0800 float64   `json:"time_allocation_0800"`
	TimeAllocation0815 float64   `json:"time_allocation_0815"`
	TimeAllocation0830 float64   `json:"time_allocation_0830"`
	TimeAllocation0845 float64   `json:"time_allocation_0845"`
	TimeAllocation0900 float64   `json:"time_allocation_0900"`
	TimeAllocation0915 float64   `json:"time_allocation_0915"`
	TimeAllocation0930 float64   `json:"time_allocation_0930"`
	TimeAllocation0945 float64   `json:"time_allocation_0945"`
	TimeAllocation1000 float64   `json:"time_allocation_1000"`
	TimeAllocation1015 float64   `json:"time_allocation_1015"`
	TimeAllocation1030 float64   `json:"time_allocation_1030"`
	TimeAllocation1045 float64   `json:"time_allocation_1045"`
	TimeAllocation1100 float64   `json:"time_allocation_1100"`
	TimeAllocation1115 float64   `json:"time_allocation_1115"`
	TimeAllocation1130 float64   `json:"time_allocation_1130"`
	TimeAllocation1145 float64   `json:"time_allocation_1145"`
	TimeAllocation1200 float64   `json:"time_allocation_1200"`
	TimeAllocation1215 float64   `json:"time_allocation_1215"`
	TimeAllocation1230 float64   `json:"time_allocation_1230"`
	TimeAllocation1245 float64   `json:"time_allocation_1245"`
	TimeAllocation1300 float64   `json:"time_allocation_1300"`
	TimeAllocation1315 float64   `json:"time_allocation_1315"`
	TimeAllocation1330 float64   `json:"time_allocation_1330"`
	TimeAllocation1345 float64   `json:"time_allocation_1345"`
	TimeAllocation1400 float64   `json:"time_allocation_1400"`
	TimeAllocation1415 float64   `json:"time_allocation_1415"`
	TimeAllocation1430 float64   `json:"time_allocation_1430"`
	TimeAllocation1445 float64   `json:"time_allocation_1445"`
	TimeAllocation1500 float64   `json:"time_allocation_1500"`
	TimeAllocation1515 float64   `json:"time_allocation_1515"`
	TimeAllocation1530 float64   `json:"time_allocation_1530"`
	TimeAllocation1545 float64   `json:"time_allocation_1545"`
	TimeAllocation1600 float64   `json:"time_allocation_1600"`
	TimeAllocation1615 float64   `json:"time_allocation_1615"`
	TimeAllocation1630 float64   `json:"time_allocation_1630"`
	TimeAllocation1645 float64   `json:"time_allocation_1645"`
	TimeAllocation1700 float64   `json:"time_allocation_1700"`
	TimeAllocation1715 float64   `json:"time_allocation_1715"`
	TimeAllocation1730 float64   `json:"time_allocation_1730"`
	TimeAllocation1745 float64   `json:"time_allocation_1745"`
	TimeAllocation1800 float64   `json:"time_allocation_1800"`
	TimeAllocation1815 float64   `json:"time_allocation_1815"`
	TimeAllocation1830 float64   `json:"time_allocation_1830"`
	TimeAllocation1845 float64   `json:"time_allocation_1845"`
	TimeAllocation1900 float64   `json:"time_allocation_1900"`
	TimeAllocation1915 float64   `json:"time_allocation_1915"`
	TimeAllocation1930 float64   `json:"time_allocation_1930"`
	TimeAllocation1945 float64   `json:"time_allocation_1945"`
	TimeAllocation2000 float64   `json:"time_allocation_2000"`
	TimeAllocation2015 float64   `json:"time_allocation_2015"`
	TimeAllocation2030 float64   `json:"time_allocation_2030"`
	TimeAllocation2045 float64   `json:"time_allocation_2045"`
	TimeAllocation2100 float64   `json:"time_allocation_2100"`
}

type ShiftTimes struct {
	StartTime float64
	EndTime   float64
	ShiftCode string
}

type WorkOrderAllocationHeaderResult struct {
	TotalTechnicianTime     float64 `json:"total_technician_time"`
	UsedTechnicianTime      float64 `json:"used_technician_time"`
	AvailableTechnicianTime float64 `json:"available_technician_time"`
	UnallocatedOperation    int     `json:"unallocated_operation"`
	AutoReleasedOperation   int     `json:"auto_released"`
	BookAllocatedTime       float64 `json:"book_allocated_time"`
}

type WorkOrderAllocationAssignTechnicianRequest struct {
	CompanyId    int       `json:"company_id" parent_entity:"trx_assign_technician" main_table:"trx_assign_technician"`
	TechnicianId int       `json:"technician_id" parent_entity:"trx_assign_technician"`
	ShiftCode    string    `json:"shift_code" parent_entity:"trx_assign_technician"`
	ForemanId    int       `json:"foreman_id" parent_entity:"trx_assign_technician"`
	ServiceDate  time.Time `json:"service_date" parent_entity:"trx_assign_technician"`
	TechnicianNo int       `json:"technician_no" parent_entity:"trx_assign_technician"`
}

type WorkOrderAllocationAssignTechnicianResponse struct {
	AssignTechnicianId int       `json:"assign_technician_id"`
	CompanyId          int       `json:"company_id"`
	CompanyName        string    `json:"company_name"`
	TechnicianId       int       `json:"technician_id"`
	TechnicianName     string    `json:"technician_name"`
	TechnicianNo       int       `json:"technician_no"`
	ShiftCode          string    `json:"shift_code"`
	ForemanId          int       `json:"foreman_id"`
	ForemanName        string    `json:"foreman_name"`
	ServiceDate        time.Time `json:"service_date"`
	Attendance         bool      `json:"attendance"`
}
