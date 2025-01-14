package transactionworkshopentities

const TableNameAtpmClaimVehicleAttachmentType = "trx_atpm_claim_vehicle_attachment_type"

type AtpmClaimVehicleAttachmentType struct {
	AttachmentTypeId          int    `gorm:"column:attachment_type_id;size:30;primaryKey" json:"attachment_type_id"`
	AttachmentCode            string `gorm:"column:attachment_code;size:30" json:"attachment_code"`
	AttachmentTypeCode        string `gorm:"column:attachment_type_code;size:30" json:"attachment_type_code"`
	AttachmentTypeDescription string `gorm:"column:attachment_type_description;size:255" json:"attachment_type_description"`
	ClaimType                 string `gorm:"column:claim_type;size:30" json:"claim_type"`
	Mandatory                 bool   `gorm:"column:mandatory" json:"mandatory"`
	MaxFileSize               int    `gorm:"column:max_file_size;size:30" json:"max_file_size"`
	IsActive                  bool   `gorm:"column:is_active" json:"is_active"`
}

func (*AtpmClaimVehicleAttachmentType) TableName() string {
	return TableNameAtpmClaimVehicleAttachmentType
}
