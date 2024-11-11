package generalservicepayloads

type Address struct {
	AddressStreet1 string `json:"address_street_1"`
	AddressStreet2 string `json:"address_street_2"`
	AddressStreet3 string `json:"address_street_3"`
	VillageId      int    `json:"village_id"`
}
type CompanyAccessData struct {
	IsActive                 int    `json:"is_active"`
	DealerRepresentativeId   string `json:"dealer_representative_id"`
	DealerRepresentativeCode string `json:"dealer_representative_code"`
	DealerRepresentativeName bool   `json:"dealer_representative_name"`
	WarehouseId              int    `json:"warehouse_id"`
}
type CompanyAccess struct {
	Page      int `json:"page"`
	PageLimit int `json:"page_limit"`
	Npages    int `json:"npages"`
	Nrows     int `json:"nrows"`
	Data      CompanyAccessData
}
type CompanyOutlet struct {
	Page      int           `json:"page"`
	PageLimit int           `json:"page_limit"`
	Npages    int           `json:"npages"`
	Nrows     int           `json:"nrows"`
	Data      []interface{} `json:"data"`
}
type BackAccount struct {
	Page      int           `json:"page"`
	PageLimit int           `json:"page_limit"`
	Npages    int           `json:"npages"`
	Nrows     int           `json:"nrows"`
	Data      []interface{} `json:"data"`
}
type EmployeeMasterResponses struct {
	IsActive          bool          `json:"is_active"`
	UserEmployeeId    int           `json:"user_employee_id"`
	UserId            int           `json:"user_id"`
	EmployeeName      string        `json:"employee_name"`
	EmployeeNickname  string        `json:"employee_nickname"`
	IdTypeId          int           `json:"id_type_id"`
	IdNumber          string        `json:"id_number"`
	CompanyId         int           `json:"company_id"`
	JobTitleId        int           `json:"job_title_id"`
	JobPositionId     int           `json:"job_position_id"`
	DivisionId        int           `json:"division_id"`
	CostCenterId      int           `json:"cost_center_id"`
	ProfitCenterId    int           `json:"profit_center_id"`
	AddressId         int           `json:"address_id"`
	Address           Address       `json:"address"`
	OfficePhoneNumber interface{}   `json:"office_phone_number"`
	HomePhoneNumber   string        `json:"home_phone_number"`
	MobilePhone       string        `json:"mobile_phone"`
	EmailAddress      string        `json:"email_address"`
	StartDate         string        `json:"start_date"`
	TerminationDate   string        `json:"termination_date"`
	GenderId          int           `json:"gender_id"`
	DateOfBirth       string        `json:"date_of_birth"`
	CityOfBirthId     int           `json:"city_of_birth_id"`
	MaritalStatusId   int           `json:"marital_status_id"`
	NumberOfChildren  int           `json:"number_of_children"`
	CitizenshipId     int           `json:"citizenship_id"`
	LastEducationId   int           `json:"last_education_id"`
	LastEmployment    string        `json:"last_employment"`
	FactorX           int           `json:"factor_x"`
	SkillLevelId      int           `json:"skill_level_id"`
	CompanyAccess     CompanyAccess `json:"company_access"`
	CompanyOutlet     CompanyOutlet `json:"company_outlet"`
	BankAccount       BackAccount   `json:"bank_account"`
}
