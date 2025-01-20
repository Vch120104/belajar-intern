package transactionworkshopentities

const TableNameAtpmMapping = "trx_atpm_mapping"

type AtpmMapping struct {
	AtpmMappingId   int    `gorm:"column:atpm_mapping_id;size:30;primaryKey" json:"atpm_mapping_id"`
	IsActive        bool   `gorm:"column:is_active" json:"is_active"`
	TableType       string `gorm:"column:table_type;size:30" json:"table_type"`
	Map0Id          int    `gorm:"column:map_0_id;size:30" json:"map_0_id"`
	Map0Code        string `gorm:"column:map_0_code;size:30" json:"map_0_code"`
	Map0Description string `gorm:"column:map_0_description;size:30" json:"map_0_description"`
	Map1Id          int    `gorm:"column:map_1_id;size:30" json:"map_1_id"`
	Map1Code        string `gorm:"column:map_1_code;size:30" json:"map_1_code"`
	Map1Description string `gorm:"column:map_1_description;size:30" json:"map_1_description"`
}

func (*AtpmMapping) TableName() string {
	return TableNameAtpmMapping
}
