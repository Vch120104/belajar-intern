package transactionworkshopentities

const TableNameAtpmClaimVehicleAttachment = "trx_atpm_claim_vehicle_attachment"

type AtpmClaimVehicleAttachment struct {
	ClaimSystemNumber int    `gorm:"column:claim_system_number;size:30;primaryKey" json:"claim_system_number"`
	CompanyId         int    `gorm:"column:company_id;size:30" json:"company_id"`
	AttachmentType    string `gorm:"column:attachment_type;size:30" json:"attachment_type"`
	FileName          string `gorm:"column:file_name;size:255" json:"file_name"`
	FileType          string `gorm:"column:file_type;size:255" json:"file_type"`
	FileSize          int    `gorm:"column:file_size;size:30" json:"file_size"`
	Extension         string `gorm:"column:extension;size:255" json:"extension"`
	Attachment        []byte `gorm:"column:attachment" json:"attachment"`
}

func (*AtpmClaimVehicleAttachment) TableName() string {
	return TableNameAtpmClaimVehicleAttachment
}
