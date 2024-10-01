package transactionjpcbpayloads

import "time"

type OutstandingJobAllocationGetAllPayload struct {
	ReferenceDocumentType   *string  `json:"reference_document_type"`
	ReferenceSystemNumber   *int     `json:"reference_system_number"`
	ReferenceDocumentNumber *string  `json:"reference_document_number"`
	VehicleId               *int     `json:"vehicle_id"`
	LeaveCar                *string  `json:"leave_car"`
	IsFromBooking           *string  `json:"is_from_booking"`
	OperationDescription    *string  `json:"operation_description"`
	ReferenceFRT            *float64 `json:"reference_frt"`
	PromiseTime             *string  `json:"promise_time"`
	LastReleaseBy           *int     `json:"last_release_by"`
	Remark                  *string  `json:"remark"`
}

type OutstandingJobAllocationGetAllResponse struct {
	ReferenceDocumentType   *string  `json:"reference_document_type"`
	ReferenceSystemNumber   *int     `json:"reference_system_number"`
	ReferenceDocumentNumber *string  `json:"reference_document_number"`
	TNKB                    *string  `json:"tnkb"`
	LeaveCar                *string  `json:"leave_car"`
	IsFromBooking           *string  `json:"is_from_booking"`
	OperationDescription    *string  `json:"operation_description"`
	ReferenceFRT            *float64 `json:"reference_frt"`
	PromiseTime             *string  `json:"promise_time"`
	LastReleaseBy           *string  `json:"last_release_by"`
	Remark                  *string  `json:"remark"`
}

type OutstandingJAProfitCenterPayload struct {
	ProfitCenterId int `json:"profit_center_id"`
}

type OutstandingJAApprovalStatusPayload struct {
	ApprovalStatusId string `json:"approval_status_id"`
}

type OutstandingJADocumentStatusPayload struct {
	DocumentStatusId int `json:"document_status_id"`
}

type OutstandingJALineTypePayload struct {
	LineTypeId int `json:"line_type_id"`
}

type OutstandingJAWorkOrderStatusPayload struct {
	WorkOrderStatusId int `json:"work_order_status_id"`
}

type OutstandingJAServiceStatusPayload struct {
	ServiceStatusId          int    `json:"service_status_id"`
	ServiceStatusDescription string `json:"service_status_description"`
}

type OutstandingJAEmployeePayload struct {
	UserEmployeeId int     `json:"user_employee_id"`
	EmployeeName   string  `json:"employee_name"`
	FactorX        float64 `json:"factor_x"`
}

type OutstandingJAVehicleMasterByTnkbPayload struct {
	VehicleId int `json:"vehicle_id"`
}

type OutstandingJobAllocationGetByTypeIdResponse struct {
	ReferenceSystemNumber   *int     `json:"reference_system_number"`
	ReferenceDocumentNumber *string  `json:"reference_document_number"`
	OperationDescription    *string  `json:"operation_description"`
	ReferenceFRT            *float64 `json:"reference_frt"`
	JobProgress             *int     `json:"job_progress"`
}

type OutstandingJobAllocationSaveRequest struct {
	CompanyId      int       `json:"company_id"`
	UserEmployeeId int       `json:"user_employee_id"`
	IsExpress      bool      `json:"is_express"`
	OperationId    int       `json:"operation_id"`
	ServiceDate    time.Time `json:"service_date"`
	SequenceNumber int       `json:"sequence_number"`
}

type OutstandingJobAllocationReferencePayload struct {
	ReferenceDocumentType   *string    `json:"reference_document_type"`
	ReferenceSystemNumber   *int       `json:"reference_system_number"`
	ReferenceDocumentNumber *string    `json:"reference_document_number"`
	ReferenceDocumentDate   *time.Time `json:"reference_document_date"`
	CostProfitCenterId      *int       `json:"cost_profit_center_id"`
	VehicleId               *int       `json:"vehicle_id"`
	Line                    *int       `json:"line"`
	FrtQuantity             *float64   `json:"frt_quantity"`
	WorkOrderDate           *time.Time `json:"work_order_date"`
	QualityControlExtraFrt  *float64   `json:"quality_control_extra_frt"`
	ReorderNumber           *float64   `json:"reorder_number"`
	BookingServiceTime      *float64   `json:"booking_service_time"`
}

type OutstandingJAVehicleMaster struct {
	VehicleId        int `json:"vehicle_id"`
	VehicleBrandId   int `json:"vehicle_brand_id"`
	VehicleModelId   int `json:"vehicle_model_id"`
	VehicleVariantId int `json:"vehicle_variant_id"`
}

type OutstandingJAVehicleStnkPayload struct {
	VehicleRegistrationCertificateTnkb string `json:"vehicle_registration_certificate_tnkb"`
}

type OutstandingJAVehicleByIdPayload struct {
	Master OutstandingJAVehicleMaster      `json:"master"`
	Stnk   OutstandingJAVehicleStnkPayload `json:"stnk"`
}

type OutstandingJAEmployeeGroupPayload struct {
	EmployeeGroupLeaderId   int    `json:"employee_group_leader_id"`
	EmployeeGroupLeaderName string `json:"employee_group_leader_name"`
	EmployeeGroupMemberId   int    `json:"employee_group_member_id"`
}

type OutstandingJobAllocationUpdateRequest struct {
	TechAllocSystemNumber int       `json:"technician_allocation_system_number"`
	CompanyId             int       `json:"company_id" validate:"required"`
	OriginalTechnicianId  int       `json:"ori_technician_id" validate:"required"`
	TechnicianId          int       `json:"technician_id" validate:"required"`
	OriSequenceNumber     int       `json:"ori_sequence_number" validate:"required"`
	SequenceNumber        int       `json:"sequence_number" validate:"required"`
	ServiceDate           time.Time `json:"service_date" validate:"required"`
}

type OutstandingJACompanyRefPayload struct {
	TimeDifference int `json:"time_difference"`
}

type OutstandingJobAllocationFetchFRTPayload struct {
	QualityControlExtraFrt float64 `json:"quality_control_extra_frt"`
	ReOrder                bool    `json:"re_order"`
	IsExpress              bool    `json:"is_express"`
	TechAllocTotalTime     float64 `json:"tech_alloc_total_time"`
	Frt                    float64 `json:"frt"`
}

type OutstandingJobAllocationSourceTargetPayload struct {
	TechnicianAllocationSystemNumber int    `json:"technician_allocation_system_number"`
	ServiceStatusId                  int    `json:"service_status_id"`
	ServiceStatusDescription         string `json:"service_status_description"`
	VehicleId                        int    `json:"vehicle_id"`
	Tnkb                             string `json:"tnkb"`
	TechnicianId                     int    `json:"technician_id"`
	TechnicianName                   string `json:"technician_name"`
	WorkOrderSystemNumber            int    `json:"work_order_system_number"`
	WorkOrderDocumentNumber          string `json:"work_order_document_number"`
	OperationCode                    string `json:"operation_code"`
}

type OutstandingJobAllocationInsertServiceLogPayload struct {
	CompanyId                        int       `json:"company_id"`
	TechnicianAllocationSystemNumber int       `json:"technician_allocation_system_number"`
	TechnicianAllocationLine         int       `json:"technician_allocation_line"`
	WorkOrderSystemNumber            int       `json:"work_order_system_number"`
	WorkOrderDocumentNumber          string    `json:"work_order_document_number"`
	WorkOrderDate                    time.Time `json:"work_order_date"`
	OperationCode                    string    `json:"operation_code"`
	TechnicianId                     int       `json:"technician_id"`
	ShiftCode                        string    `json:"shift_code"`
	FrtQuantity                      float64   `json:"frt_quantity"`
	TechAllocLastStartDate           time.Time `json:"tech_alloc_last_start_date"`
	TechAllocLastEndDate             time.Time `json:"tech_alloc_last_end_date"`
	TechAllocLastStartTime           float64   `json:"tech_alloc_last_start_time"`
	SequenceNumber                   int       `json:"sequence_number"`
}

type OutstandingJobAllocationUpdateResponse struct {
	SourceFirstTechAllocSystenNumber int `json:"source_first_tech_alloc_system_number"`
	TechAllocSystemNumber            int `json:"tech_alloc_system_number"`
}
