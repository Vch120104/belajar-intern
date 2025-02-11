package transactionworkshopentities

import "time"

const TableNamePrintGatePass = "trx_gate_pass"

type PrintGatePass struct {
	CompanyId              int       `gorm:"column:company_id;size:30;" json:"company_id"`
	GatePassSystemNumber   int       `gorm:"column:gate_pass_system_number;primaryKey;size:30;" json:"gate_pass_system_number"`
	GatePassDocumentNumber string    `gorm:"column:gate_pass_document_number;size:100;" json:"gate_pass_document_number"`
	GatePassDate           time.Time `gorm:"column:gate_pass_date;" json:"gate_pass_date"`
	VehicleId              int       `gorm:"column:vehicle_id;size:30;" json:"vehicle_id"`
	BrandId                int       `gorm:"column:brand_id;size:30;" json:"brand_id"`
	ModelId                int       `gorm:"column:model_id;size:30;" json:"model_id"`
	VariantId              int       `gorm:"column:variant_id;size:30;" json:"variant_id"`
	ColourId               int       `gorm:"column:colour_id;size:30;" json:"colour_id"`
	BpkSystemNumber        int       `gorm:"column:bpk_system_number;size:30;" json:"bpk_system_number"`
	TrfOutSystemNumber     int       `gorm:"column:trf_out_system_number;size:30;" json:"trf_out_system_number"`
	TrfOutDocumentNumber   string    `gorm:"column:trf_out_document_number;size:100;" json:"trf_out_document_number"`
	TrfOutDocumentDate     time.Time `gorm:"column:trf_out_document_date;" json:"trf_out_document_date"`
	TrfOutLineNumber       int       `gorm:"column:trf_out_line_number;size:30;" json:"trf_out_line_number"`
	CustomerId             int       `gorm:"column:customer_id;size:30;" json:"customer_id"`
	WarehouseId            int       `gorm:"column:warehouse_id;size:30;" json:"warehouse_id"`
	DeliveryName           string    `gorm:"column:delivery_name;size:100;" json:"delivery_name"`
	DeliveryAddress        string    `gorm:"column:delivery_address;size:256;" json:"delivery_address"`
	VillageId              int       `gorm:"column:village_id;size:30;" json:"village_id"`
	DistrictId             int       `gorm:"column:district_id;size:30;" json:"district_id"`
	ProvinceId             int       `gorm:"column:province_id;size:30;" json:"province_id"`
	CityId                 int       `gorm:"column:city_id;size:30;" json:"city_id"`
	LastPrintBy            string    `gorm:"column:last_print_by;size:50;" json:"last_print_by"`
	PrintingNumber         int       `gorm:"column:printing_number;size:30;" json:"printing_number"`
}

func (*PrintGatePass) TableName() string {
	return TableNamePrintGatePass
}
