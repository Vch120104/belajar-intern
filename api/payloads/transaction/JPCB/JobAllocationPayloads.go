package transactionjpcbpayloads

type GetAllJobAllocationPayload struct {
	TechnicianAllocationSystemNumber *int     `json:"technician_allocation_system_number"`
	TechnicianId                     *int     `json:"technician_id"`
	ServiceStatusId                  *int     `json:"service_status_id"`
	SequenceNumber                   *int     `json:"sequence_number"`
	ReferenceDocumentType            *string  `json:"reference_document_type"`
	ReferenceDocumentNumber          *string  `json:"reference_document_number"`
	VehicleId                        *int     `json:"vehicle_id"`
	Operation                        *string  `json:"operation"`
	Frt                              *float64 `json:"frt"`
	FactorX                          *float64 `json:"factor_x"`
	FrtJPCB                          *float64 `json:"frt_jpcb"`
	TechAllocLastStartTime           *float64 `json:"tech_alloc_last_start_time"`
	TechAllocLastEndTime             *float64 `json:"tech_alloc_last_end_time"`
	IsExpress                        *bool    `json:"is_express"`
}

type ItemGroupPayload struct {
	ItemGroupId int `json:"item_group_id"`
}

type VehicleMaster struct {
	VehicleId int `json:"vehicle_id"`
}

type VehicleStnkPayload struct {
	VehicleRegistrationCertificateTnkb string `json:"vehicle_registration_certificate_tnkb"`
}

type VehiclePayload struct {
	Master VehicleMaster      `json:"master"`
	Stnk   VehicleStnkPayload `json:"stnk"`
}

type UserDetailsPayload struct {
	UserEmployeeId int    `json:"user_employee_id"`
	EmployeeName   string `json:"employee_name"`
}

type ServiceStatusPayload struct {
	ServiceStatusId          int    `json:"service_status_id"`
	ServiceStatusDescription string `json:"service_status_description"`
}

type GetAllJobAllocationResponse struct {
	TechnicianAllocationSystemNumber *int     `json:"technician_allocation_system_number"`
	TechnicianName                   *string  `json:"technician_name"`
	ServiceStatus                    *string  `json:"service_status"`
	SequenceNumber                   *int     `json:"sequence_number"`
	ReferenceDocumentType            *string  `json:"reference_document_type"`
	ReferenceDocumentNumber          *string  `json:"reference_document_number"`
	VehicleTNBK                      *string  `json:"vehicle_tnbk"`
	Operation                        *string  `json:"operation"`
	Frt                              *float64 `json:"frt"`
	FactorX                          *float64 `json:"factor_x"`
	FrtJPCB                          *float64 `json:"frt_jpcb"`
	TechAllocLastStartTime           *float64 `json:"tech_alloc_last_start_time"`
	TechAllocLastEndTime             *float64 `json:"tech_alloc_last_end_time"`
	IsExpress                        *bool    `json:"is_express"`
}

type GetJobAllocationByIdPayload struct {
	TechnicianAllocationSystemNumber *int     `json:"technician_allocation_system_number"`
	CompanyId                        *int     `json:"company_id"`
	TechnicianId                     *int     `json:"technician_id"`
	SequenceNumber                   *int     `json:"sequence_number"`
	ReferenceDocumentNumber          *string  `json:"reference_document_number"`
	Operation                        *string  `json:"operation"`
	Frt                              *float64 `json:"frt"`
	FactorX                          *float64 `json:"factor_x"`
	WorkOrderSystemNumber            *int     `json:"work_order_system_number"`
}

type GetProgressResponse struct {
	Progress float64 `json:"progress"`
}

type GetJobAllocationByIdResponse struct {
	TechnicianAllocationSystemNumber *int     `json:"technician_allocation_system_number"`
	Operation                        *string  `json:"operation"`
	Frt                              *float64 `json:"frt"`
	FactorX                          *float64 `json:"factor_x"`
	TechnicianName                   *string  `json:"technician_name"`
	SequenceNumber                   *int     `json:"sequence_number"`
	Progress                         *float64 `json:"progress"`
}

type JobAllocationUpdateRequest struct {
	TechnicianId   int `json:"technician_id"`
	SequenceNumber int `json:"sequence_number"`
}
