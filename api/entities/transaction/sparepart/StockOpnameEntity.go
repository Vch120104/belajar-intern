package transactionsparepartentities

import (
	"database/sql"
	"time"
)

type AtStockOpname0 struct {
	RecordStatus      string       `gorm:"column:RECORD_STATUS" json:"record_status"`
	CompanyCode       float64      `gorm:"column:COMPANY_CODE" json:"company_code"`
	StockOpnameSysNo  float64      `gorm:"column:STOCK_OPNAME_SYS_NO;autoIncrement" json:"stock_opname_sys_no"`
	StockOpnameDocNo  string       `gorm:"column:STOCK_OPNAME_DOC_NO" json:"stock_opname_doc_no"`
	StockOpnameStatus string       `gorm:"column:STOCK_OPNAME_STATUS" json:"stock_opname_status"`
	WhsGroup          string       `gorm:"column:WHS_GROUP" json:"whs_group"`
	WhsCode           string       `gorm:"column:WHS_CODE" json:"whs_code"`
	LocRangeFrom      string       `gorm:"column:LOC_RANGE_FROM" json:"loc_range_from"`
	LocRangeTo        string       `gorm:"column:LOC_RANGE_TO" json:"loc_range_to"`
	ProfitCenter      string       `gorm:"column:PROFIT_CENTER" json:"profit_center"`
	TrxType           string       `gorm:"column:TRX_TYPE" json:"trx_type"`
	EventNo           string       `gorm:"column:EVENT_NO" json:"event_no"`
	ShowDetail        bool         `gorm:"column:SHOW_DETAIL" json:"show_detail"`
	Pic               string       `gorm:"column:PIC" json:"pic"`
	Remark            string       `gorm:"column:REMARK" json:"remark"`
	ItemGroup         string       `gorm:"column:ITEM_GROUP" json:"item_group"`
	ExecDateFrom      time.Time    `gorm:"column:EXEC_DATE_FROM;type:dateTime" json:"exec_date_from"`
	ExecDateTo        time.Time    `gorm:"column:EXEC_DATE_TO;type:dateTime" json:"exec_date_to"`
	BrokenWhsCode     string       `gorm:"column:BROKEN_WHS_CODE" json:"broken_whs_code"`
	AdjustDate        sql.NullTime `gorm:"column:ADJUST_DATE" json:"adjust_date"`
	ApprovalStatus    string       `gorm:"column:APPROVAL_STATUS" json:"approval_status"`
	ApprovalReqBy     string       `gorm:"column:APPROVAL_REQ_BY" json:"approval_req_by"`
	ApprovalReqDate   sql.NullTime `gorm:"column:APPROVAL_REQ_DATE" json:"approval_req_date"`
	ApprovalBy        string       `gorm:"column:APPROVAL_BY" json:"approval_by"`
	ApprovalDate      sql.NullTime `gorm:"column:APPROVAL_DATE" json:"approval_date"`
	ChangeNo          float64      `gorm:"column:CHANGE_NO" json:"change_no"`
	CreationUserId    string       `gorm:"column:CREATION_USER_ID" json:"creation_user_id"`
	CreationDatetime  time.Time    `gorm:"column:CREATION_DATETIME" json:"creation_datetime"`
	ChangeUserId      string       `gorm:"column:CHANGE_USER_ID" json:"change_user_id"`
	ChangeDatetime    time.Time    `gorm:"column:CHANGE_DATETIME" json:"change_datetime"`
	TotalAdjCost      float64      `gorm:"column:TOTAL_ADJ_COST" json:"total_adj_cost"`
	IncludeZeroOnhand sql.NullBool `gorm:"column:Include_Zero_Onhand" json:"include_zero_onhand"`
}

func (*AtStockOpname0) TableName() string {
	return "atStockOpname0"
}

type GmLoc2 struct {
	RecordStatus     string    `gorm:"column: RECORD_STATUS" json:"record_status"`
	CompanyCode      float64   `gorm:"column: COMPANY_CODE" json:"company_code"`
	WarehouseCode    string    `gorm:"column: WAREHOUSE_CODE" json:"warehouse_code"`
	LocationCode     string    `gorm:"column: LOCATION_CODE" json:"location_code"`
	WarehouseGroup   string    `gorm:"column: WAREHOUSE_GROUP" json:"warehouse_group"`
	LocationName     string    `gorm:"column: LOCATION_NAME" json:"location_name"`
	Description      string    `gorm:"column: DESCRIPTION" json:"description"`
	PickSequence     int16     `gorm:"column: PICK_SEQUENCE" json:"pick_sequence"`
	CapacityInM3     float64   `gorm:"column: CAPACITY_IN_M3" json:"capacity_in_m3"`
	ChangeNo         float64   `gorm:"column: CHANGE_NO" json:"change_no"`
	CreationUserId   string    `gorm:"column: CREATION_USER_ID" json:"creation_user_id"`
	CreationDatetime time.Time `gorm:"column: CREATION_DATETIME" json:"creation_datetime"`
	ChangeUserId     string    `gorm:"column: CHANGE_USER_ID" json:"change_user_id"`
	ChangeDatetime   time.Time `gorm:"column: CHANGE_DATETIME" json:"change_datetime"`
}

func (*GmLoc2) TableName() string {
	return "gmLoc2"
}

type GmEmp struct {
	RecordStatus         string    `gorm:"column: RECORD_STATUS" json:"record_status"`
	EmployeeNo           string    `gorm:"column: EMPLOYEE_NO" json:"employee_no"`
	CompanyCode          float64   `gorm:"column: COMPANY_CODE" json:"company_code"`
	EmployeeName         string    `gorm:"column: EMPLOYEE_NAME" json:"employee_name"`
	EmpNickName          string    `gorm:"column: EMP_NICK_NAME" json:"emp_nick_name"`
	IdType               string    `gorm:"column: ID_TYPE" json:"id_type"`
	IdNo                 string    `gorm:"column: ID_NO" json:"id_no"`
	JobTitle             string    `gorm:"column: JOB_TITLE" json:"job_title"`
	JobPosition          string    `gorm:"column: JOB_POSITION" json:"job_position"`
	CostCenter           string    `gorm:"column: COST_CENTER" json:"cost_center"`
	ProfitCenter         string    `gorm:"column: PROFIT_CENTER" json:"profit_center"`
	SkillLvl             string    `gorm:"column: SKILL_LVL" json:"skill_lvl"`
	AccountNo            string    `gorm:"column: ACCOUNT_NO" json:"account_no"`
	AccountName          string    `gorm:"column: ACCOUNT_NAME" json:"account_name"`
	BankCode             string    `gorm:"column: BANK_CODE" json:"bank_code"`
	BankBranchCode       string    `gorm:"column: BANK_BRANCH_CODE" json:"bank_branch_code"`
	Department           string    `gorm:"column: DEPARTMENT" json:"department"`
	Manager              string    `gorm:"column: MANAGER" json:"manager"`
	OfficePhoneNo        string    `gorm:"column: OFFICE_PHONE_NO" json:"office_phone_no"`
	EmailAddress         string    `gorm:"column: EMAIL_ADDRESS" json:"email_address"`
	HomeAddress1         string    `gorm:"column: HOME_ADDRESS_1" json:"home_address_1"`
	HomeAddress2         string    `gorm:"column: HOME_ADDRESS_2" json:"home_address_2"`
	HomeAddress3         string    `gorm:"column: HOME_ADDRESS_3" json:"home_address_3"`
	HomeVillageCode      string    `gorm:"column: HOME_VILLAGE_CODE" json:"home_village_code"`
	HomeSubdistrictCode  string    `gorm:"column: HOME_SUBDISTRICT_CODE" json:"home_subdistrict_code"`
	HomeMunicipalityCode string    `gorm:"column: HOME_MUNICIPALITY_CODE" string:"home_municipality_code"`
	HomeProvinceCode     string    `gorm:"column: HOME_PROVINCE_CODE" json:"home_province_code"`
	HomeCityCode         string    `gorm:"column: HOME_CITY_CODE" json:"home_city_code"`
	HomeZipCode          string    `gorm:"column: HOME_ZIP_CODE" json:"home_zip_code"`
	HomePhoneNo          string    `gorm:"column: HOME_PHONE_NO" json:"home_phone_no"`
	MobilePhone          string    `gorm:"column: MOBILE_PHONE" json:"mobile_phone"`
	Gender               string    `gorm:"column: GENDER" json:"gender"`
	DateOfBirth          time.Time `gorm:"column: DATE_OF_BIRTH" json:"date_of_birth"`
	CityOfBirth          string    `gorm:"column: CITY_OF_BIRTH" json:"city_of_birth"`
	MaritalStatus        string    `gorm:"column: MARITAL_STATUS" json:"marital_status"`
	NoOFChildern         float64   `gorm:"column: NO_OF_CHILDREN" json:"no_of_children"`
	Citizenship          string    `gorm:"column: CITIZENSHIP" json:"citizenship"`
	StartDate            time.Time `gorm:"column: START_DATE" json:"start_date"`
	TerminationDate      time.Time `gorm:"column: TERMINATION_DATE" json:"termination_date"`
	LastEducation        string    `gorm:"column: LAST_EDUCATION" json:"last_education"`
	LastEmployment       string    `gorm:"column: LAST_EMPLOYMENT" json:"last_employment"`
	BrGroupCode          string    `gorm:"column: BR_GROUP_CODE" json:"br_group_code"`
	KeyPass              string    `gorm:"column: KEY_PASS" json:"key_pass"`
	HassPass             string    `gorm:"column: HASS_PASS" json:"hass_pass"`
	ChangeNo             float64   `gorm:"column: CHANGE_NO" json:"change_no"`
	ChangeUserId         string    `gorm:"column: CHANGE_USER_ID" json:"change_user_id"`
	ChangeDatetime       time.Time `gorm:"column: CHANGE_DATETIME" json:"change_datetime"`
	CreationUserId       string    `gorm:"column: CREATION_USER_ID" json:"creation_user_id"`
	CreationDatetime     time.Time `gorm:"column: CREATION_DATETIME" json:"creation_datetime"`
	FactorX              float64   `gorm:"column: FACTOR_X" json:"factor_x"`
	SkillLevelCode       string    `gorm:"column: SKILL_LEVEL_CODE" json:"skill_level_code"`
	TokenNotification    string    `gorm:"column: TOKEN_NOTIFICATION" json:"token_notification"`
	TokenInstanceid      string    `gorm:"column: TOKEN_INSTANCEID" json:"token_instanceid"`
	ForceExpire          string    `gorm:"column: FORCE_EXPIRE" json:"force_expire"`
}

func (ge *GmEmp) TableName() string {
	return "gmEmp"
}

type AtStockOpname1 struct {
	RecordStatus     string    `gorm:"column: RECORD_STATUS" json:"record_status"`
	StockOpnameSysNo float64   `gorm:"column: STOCK_OPNAME_SYS_NO;autoIncrement" json:"stock_opname_sys_no"`
	StockOpnameLine  float64   `gorm:"column: STOCK_OPNAME_LINE;autoIncrement" json:"stock_opname_line"`
	LineStatus       string    `gorm:"column: LINE_STATUS" json:"line_status"`
	WhsCode          string    `gorm:"column: WHS_CODE" json:"whs_code"`
	LocCode          string    `gorm:"column: LOC_CODE" json:"loc_code"`
	ItemCode         string    `gorm:"column: ITEM_CODE" json:"item_code"`
	SysQty           float64   `gorm:"column: SYS_QTY" json:"sys_qty"`
	FoundQty         float64   `gorm:"column: FOUND_QTY" json:"found_qty"`
	BrokeQty         float64   `gorm:"column: BROKE_QTY" json:"broke_qty"`
	BrokenLocCode    string    `gorm:"column: BROKEN_LOC_CODE" json:"broken_loc_code"`
	NeedAdjustment   bool      `gorm:"column: NEED_ADJUSTMENT" json:"need_adjustment"`
	Remark           string    `gorm:"column: REMARK" json:"remark"`
	ChangeNo         float64   `gorm:"column: CHANGE_NO" json:"change_no"`
	CreationUserId   string    `gorm:"column: CREATION_USER_ID" json:"creation_user_id"`
	CreationDatetime time.Time `gorm:"column: CREATION_DATETIME" json:"creation_datetime"`
	ChangeUserId     string    `gorm:"column: CHANGE_USER_ID" json:"change_user_id"`
	ChangeDatetime   time.Time `gorm:"column: CHANGE_DATETIME" json:"change_datetime"`
	Cogs             float64   `gorm:"column: COGS" json:"cogs"`
	AdjCost          float64   `gorm:"column: ADJ_COST" json:"adj_cost"`
	AllocQty         float64   `gorm:"column: ALLOC_QTY" json:"alloc_qty"`
}

func (*AtStockOpname1) TableName() string {
	return "atStockOpname1"
}
