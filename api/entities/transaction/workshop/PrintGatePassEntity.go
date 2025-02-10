package transactionworkshopentities

const TableNamePrintGatePass = "trx_gate_pass"

type PrintGatePass struct {
	CompanyId              int    `gorm:"column:company_id;not null" json:"company_id"`
	GatePassSystemNumber   int    `gorm:"column:gate_pass_system_number;primary_key;not null" json:"gate_pass_system_number"`
	GatePassDocumentNumber string `gorm:"column:gate_pass_document_number;size:25;not null" json:"gate_pass_document_number"`
	GatePassDate           string `gorm:"column:gate_pass_date;type:datetime;not null" json:"gate_pass_date"`
	VehicleId              int    `gorm:"column:vehicle_id;not null" json:"vehicle_id"`
	VehicleBrandId         int    `gorm:"column:vehicle_brand_id;not null" json:"vehicle_brand_id"`
	ModelId                int    `gorm:"column:model_id;not null" json:"model_id"`
	VariantId              int    `gorm:"column:variant_id;not null" json:"variant_id"`
	ColourId               int    `gorm:"column:colour_id;not null" json:"colour_id"`
	BpkSystemNumber        int    `gorm:"column:bpk_system_number;not null" json:"bpk_system_number"`
	TrfOutSystemNumber     int    `gorm:"column:trf_out_system_number;not null" json:"trf_out_system_number"`
	TrfOutDocumentNumber   string `gorm:"column:trf_out_document_number;size:25;not null" json:"trf_out_document_number"`
	TrfOutDocumentDate     string `gorm:"column:trf_out_document_date;type:datetime;not null" json:"trf_out_document_date"`
	TrfOutLineNumber       int    `gorm:"column:trf_out_line_number;not null" json:"trf_out_line_number"`
	CustomerId             int    `gorm:"column:customer_id;not null" json:"customer_id"`
	WarehouseId            int    `gorm:"column:warehouse_id;not null" json:"warehouse_id"`
	DeliveryName           string `gorm:"column:delivery_name;size:100;not null" json:"delivery_name"`
	DeliveryAddress        string `gorm:"column:delivery_address;size:256;not null" json:"delivery_address"`
	VillageId              string `gorm:"column:village_id;size:10;not null" json:"village_id"`
	DistrictId             string `gorm:"column:district_id;size:10;not null" json:"district_id"`
	ProvinceId             string `gorm:"column:province_id;size:10;not null" json:"province_id"`
	CityId                 string `gorm:"column:city_id;size:5;not null" json:"city_id"`
	LastPrintBy            string `gorm:"column:last_print_by;size:50;not null" json:"last_print_by"`
	PrintingNumber         int    `gorm:"column:printing_number;not null" json:"printing_number"`
}

func (*PrintGatePass) TableName() string {
	return TableNamePrintGatePass
}
