package transactionworkshopentities

import "time"

const TableNameWorkOrderAllocation = "trx_work_order_allocation"

type WorkOrderAllocation struct {
	WorkOrderAllocationId   int       `gorm:"column:work_order_allocation_id;size:30;primaryKey" json:"work_order_allocation_id"`
	IsActive                bool      `gorm:"column:is_active;size:1" json:"is_active"`
	CompanyId               int       `gorm:"column:company_id;size:30" json:"company_id"`
	BrandId                 int       `gorm:"column:brand_id;size:30" json:"brand_id"`
	ProfitCenterId          int       `gorm:"column:profit_center_id;size:30" json:"profit_center_id"`
	TechAllocSystemNumber   int       `gorm:"column:tech_alloc_system_number;size:30" json:"tech_alloc_system_number"`
	TechnicianId            int       `gorm:"column:technician_id;size:30" json:"technician_id"`
	ForemanId               int       `gorm:"column:foreman_id;size:30" json:"foreman_id"`
	UsingGroup              bool      `gorm:"column:using_group;size:1" json:"using_group"`
	TechnicianGroup         string    `gorm:"column:technician_group;" json:"technician_group"`
	SequenceNumber          int       `gorm:"column:sequence_number;size:30" json:"sequence_number"`
	TechAllocStartDate      time.Time `gorm:"column:tech_alloc_start_date" json:"tech_alloc_start_date"`
	TechAllocEndDate        time.Time `gorm:"column:tech_alloc_end_date" json:"tech_alloc_end_date"`
	TechAllocStartTime      float64   `gorm:"column:tech_alloc_start_time" json:"tech_alloc_start_time"`
	TechAllocEndTime        float64   `gorm:"column:tech_alloc_end_time" json:"tech_alloc_end_time"`
	TechAllocTotalTime      float64   `gorm:"column:tech_alloc_total_time" json:"tech_alloc_total_time"`
	TechAllocLastStartDate  time.Time `gorm:"column:tech_alloc_last_start_date" json:"tech_alloc_last_start_date"`
	TechAllocLastEndDate    time.Time `gorm:"column:tech_alloc_last_end_date" json:"tech_alloc_last_end_date"`
	TechAllocLastStartTime  float64   `gorm:"column:tech_alloc_last_start_time" json:"tech_alloc_last_start_time"`
	TechAllocLastEndTime    float64   `gorm:"column:tech_alloc_last_end_time" json:"tech_alloc_last_end_time"`
	OperationCode           string    `gorm:"column:operation_code;" json:"operation_code"`
	ShiftCode               string    `gorm:"column:shift_code;" json:"shift_code"`
	ServActualTime          float64   `gorm:"column:serv_actual_time" json:"serv_actual_time"`
	ServPendingTime         float64   `gorm:"column:serv_pending_time" json:"serv_pending_time"`
	ServProgressTime        float64   `gorm:"column:serv_progress_time" json:"serv_progress_time"`
	ServTotalActualTime     float64   `gorm:"column:serv_total_actual_time" json:"serv_total_actual_time"`
	ServStatus              string    `gorm:"column:serv_status;" json:"serv_status"`
	BookingSystemNumber     int       `gorm:"column:booking_system_number;size:30" json:"booking_system_number"`
	BookingDocumentNumber   string    `gorm:"column:booking_document_number;" json:"booking_document_number"`
	BookingLine             int       `gorm:"column:booking_line;size:30" json:"booking_line"`
	WorkOrderSystemNumber   int       `gorm:"column:work_order_system_number;size:30" json:"work_order_system_number"`
	WorkOrderDocumentNumber string    `gorm:"column:work_order_document_number;" json:"work_order_document_number"`
	WorkOrderLine           int       `gorm:"column:work_order_line;size:30" json:"work_order_line"`
	ReOrder                 bool      `gorm:"column:re_order;size:1" json:"re_order"`
	InvoiceSystemNumber     int       `gorm:"column:invoice_system_number;size:30" json:"invoice_system_number"`
	InvoiceDocumentNumber   string    `gorm:"column:invoice_document_number;" json:"invoice_document_number"`
	ChangeNo                int       `gorm:"column:change_no;size:30" json:"change_no"`
	CreationDate            time.Time `gorm:"column:creation_date" json:"creation_date"`
	CreatedBy               int       `gorm:"column:created_by;size:30" json:"created_by"`
	ChangeDate              time.Time `gorm:"column:change_date" json:"change_date"`
	ChangedBy               int       `gorm:"column:changed_by;size:30" json:"changed_by"`
	FactorX                 float64   `gorm:"column:factor_x" json:"factor_x"`
	IsExpress               bool      `gorm:"column:is_express;size:1" json:"is_express"`
	Frt                     float64   `gorm:"column:frt" json:"frt"`
	BookingServiceTime      float64   `gorm:"column:booking_service_time" json:"booking_service_time"`
}

func (*WorkOrderAllocation) TableName() string {
	return TableNameWorkOrderAllocation
}
