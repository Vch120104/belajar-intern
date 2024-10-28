package masterentities

const ItemClaimTypeTableName = "mtr_item_claim_type"

type ItemClaimType struct {
	ItemClaimTypeId          int    `gorm:"column:item_claim_type_id;not null;primaryKey;size:30" json:"item_claim_type_id"`
	ItemClaimTypeCode        string `gorm:"column:item_claim_type_code;not null;size:10" json:"item_claim_type_code"`
	ItemClaimTypeDescription string `gorm:"column:item_claim_type_description;not null;size:10" json:"item_claim_type_description"`
}

func (*ItemClaimType) TableName() string {
	return ItemClaimTypeTableName
}
