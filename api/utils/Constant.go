package utils

import (
	"errors"
	"strconv"
	"time"
)

// WO and Service Status
var SrvStatDraft int = 1       // 0 Draft
var SrvStatStart int = 2       // 10 Start
var SrvStatPending int = 3     // 15 Pending
var SrvStatStop int = 4        // 20 Stop
var SrvStatTransfer int = 5    // 25 Transfer
var SrvStatQcPass int = 6      // 30 QC Pass
var SrvStatReOrder int = 7     // 35 Re-Order
var SrvStatAutoRelease int = 8 // 40 Auto Release

var WoStatDraft int = 1   // 0 Draft
var WoStatNew int = 2     // 10 New
var WoStatReady int = 3   // 20 Ready
var WoStatOngoing int = 4 // 30 On Going
var WoStatStop int = 5    // 40 Stop
var WoStatQC int = 6      // 50 QC
var WoStatCancel int = 7  // 60 Cancel
var WoStatClosed int = 8  // 70 Close

// Linetype Status
var LinetypePackage int = 0            // 0 Package BodyShop
var LinetypeOperation int = 1          // 1 Operation
var LinetypeSparepart int = 2          // 2 Sparepart
var LinetypeOil int = 3                // 3 Oil
var LinetypeMaterial int = 4           // 4 Material
var LinetypeFee int = 5                // 5 Fee
var LinetypeAccesories int = 6         // 6 Accesories
var LinetypeConsumableMaterial int = 7 // 7 Consumable Material
var LinetypeSublet int = 8             // 8 Sublet
var LinetypeSouvenir int = 9           // 9 Souvenir

// Transaction Type Bill Code SO WO
type TrxType struct {
	Code string
	ID   int
}

var TrxTypeWoInternal = TrxType{Code: "I", ID: 1}        // TRXTYPE_WO_INTERNAL
var TrxTypeWoNoCharge = TrxType{Code: "N", ID: 2}        // TRXTYPE_WO_NOCHARGE
var TrxTypeWoCentralize = TrxType{Code: "C", ID: 3}      // TRXTYPE_WO_CENTRALIZE
var TrxTypeWoDeCentralize = TrxType{Code: "D", ID: 4}    // TRXTYPE_WO_DECENTRALIZE
var TrxTypeWoCampaign = TrxType{Code: "G", ID: 5}        // TRXTYPE_WO_CAMPAIGN
var TrxTypeWoContractService = TrxType{Code: "S", ID: 6} // TRXTYPE_WO_CONTRACT_SERVICE
var TrxTypeWoExternal = TrxType{Code: "E", ID: 7}        // TRXTYPE_WO_EXTERNAL
var TrxTypeWoFreeService = TrxType{Code: "F", ID: 8}     // TRXTYPE_WO_FREE_SERVICE
var TrxTypeWoInsurance = TrxType{Code: "U", ID: 9}       // TRXTYPE_WO_INSURANCE
var TrxTypeWoWarranty = TrxType{Code: "W", ID: 10}       // TRXTYPE_WO_WARRANTY

var TrxTypeSoDirect = TrxType{Code: "SU01", ID: 1}       // TRXTYPE_SO_DIRECT
var TrxTypeSoChannel = TrxType{Code: "SU02", ID: 2}      // TRXTYPE_SO_CHANNEL
var TrxTypeSoGSO = TrxType{Code: "SU03", ID: 3}          // TRXTYPE_SO_GSO
var TrxTypeSoInternal = TrxType{Code: "SU05", ID: 4}     // TRXTYPE_SO_INTERNAL
var TrxTypeSoCentralize = TrxType{Code: "SU06", ID: 5}   // TRXTYPE_SO_CENTRALIZE
var TrxTypeSoDeCentralize = TrxType{Code: "SU07", ID: 6} // TRXTYPE_SO_DECENTRALIZED
var TrxTypeSoExport = TrxType{Code: "SU08", ID: 7}       // TRXTYPE_SO_EXPORT

// Job Type
type JobType struct {
	Code string
	ID   int
}

var JobTypeBodyRepair = JobType{Code: "BR", ID: 1}             // JOBTYPE_BODYREPAIR
var JobTypeCampaign = JobType{Code: "CP", ID: 2}               // JOBTYPE_CAMPAIGN
var JobTypeContractService = JobType{Code: "CS", ID: 3}        // JOBTYPE_CONTRACTSERVICE
var JobTypeFreeServiceInspection = JobType{Code: "FSI", ID: 4} // JOBTYPE_FREESERVICEINSPECTION
var JobTypeGeneralRepair = JobType{Code: "GR", ID: 5}          // JOBTYPE_GENERALREPAIR
var JobTypeJPAccessories = JobType{Code: "JPA", ID: 6}         // JOBTYPE_JPACCESSORIES
var JobTypeMarketing = JobType{Code: "M", ID: 7}               // JOBTYPE_MARKETING
var JobTypePDI = JobType{Code: "PDI", ID: 8}                   // JOBTYPE_PDI
var JobTypePeriodicalMaintenance = JobType{Code: "PM", ID: 9}  // JOBTYPE_PERIODICALMAINTENANCE
var JobTypePurchasing = JobType{Code: "P", ID: 10}             // JOBTYPE_PURCHASING
var JobTypeRobbing = JobType{Code: "RJ", ID: 11}               // JOBTYPE_ROBBING
var JobTypeTB = JobType{Code: "TB", ID: 12}                    // JOBTYPE_TRANSFERTOBODYREPAIR
var JobTypeTG = JobType{Code: "TG", ID: 13}                    // JOBTYPE_TRANSFERTOGENERALREPAIR
var JobTypeWarehouse = JobType{Code: "W", ID: 14}              // JOBTYPE_WAREHOUSE
var JobTypeWarranty = JobType{Code: "W", ID: 15}               // JOBTYPE_WARRANTY

var UomTypeService string = "S"     // UOMTYPE_SERVICE Services
var EstWoOrderTypeId int = 1        // EST_WO_ORDER_TYPE Order Type For Work Order and Estimation
var EstWoOrderType string = "E"     // EST_WO_ORDER_TYPE Order Type For Work Order and Estimation
var EstWoDiscSelectionId int = 1    // EST_WO_DISC_SELECTION Discount Selection for Estimation and WO
var EstWoDiscSelection string = "D" // EST_WO_DISC_SELECTION Discount Selection for Estimation and WO

// CarWash
var CarWashStatDraft int = 1
var CarWashStatStart int = 2
var CarWashStatStop int = 3

var CarWashPriorityHigh int = 1
var CarWashPriorityNormal int = 2

// Status

var Draft int = 1

var Revise int = 99

// Login
var LoginSuccess string = "Login Success"
var LoginFailed string = "Login Failed"

// Success
var GetDataSuccess string = "Get Data Successfully"
var CreateDataSuccess string = "Create Data Successfully"
var UpdateDataSuccess string = "Update Data Successfully"
var DeleteDataSuccess string = "Delete Data Successfully"

// Failed
var GetDataFailed string = "Get Data Failed"
var CreateDataFailed string = "Create Data Successfully"
var UpdateDataFailed string = "Update Data Failed"
var DeleteDataFailed string = "Delete Data Failed"

// Error
var CannotSendEmail string = "Cannot Send Email"
var DataExists string = "Data Already Exists"
var GetDataNotFound string = "Data Not Found"
var SomethingWrong string = "Something wrong, please contact admin"
var ServiceUnavailable string = "Service Unavailable"
var InvalidInput string = "Invalid Input"
var BadRequestError string = "Please check your input"
var JsonError string = "Please check your json input"
var SessionError string = "Session Invalid, please re-login"
var MultiLoginError string = "you are already logged in on a different device"
var PermissionError string = "You don't have permission"
var PasswordNotMatched string = "Password not matched"
var ExcelEpoch = time.Date(1899, time.December, 30, 0, 0, 0, 0, time.UTC)

// Etc
var LikeString string = "%%%s%%"

func BoolPtr(b bool) *bool {
	return &b
}
func IntPtr(i int) *int {
	return &i
}
func TimePtr(t time.Time) *time.Time {
	return &t
}

func StringPtr(str string) *string {
	return &str
}

func RemoveDuplicateIds(arr []int) []int {
	encountered := make(map[int]bool)
	result := []int{}

	for _, v := range arr {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}

	return result
}

func ExcelDateToDate(excelDate string) time.Time {
	var days, _ = strconv.ParseFloat(excelDate, 64)
	return ExcelEpoch.Add(time.Second * time.Duration(days*86400))
}

// Error
var ErrIncorrectInput = errors.New(BadRequestError)
var ErrNotFound = errors.New(GetDataNotFound)
var ErrConflict = errors.New(DataExists)
var ErrEntity = errors.New(JsonError)
var ErrInternalServerError = errors.New(SomethingWrong)
var ErrServiceUnavailable = errors.New(ServiceUnavailable)
var ErrSession = errors.New(SessionError)
var ErrMultiLogin = errors.New(MultiLoginError)
var ErrPermission = errors.New(PermissionError)
var ErrPasswordNotMatched = errors.New(PasswordNotMatched)
var ErrCannotSendEmail = errors.New(CannotSendEmail)
var ErrInvalidInput = errors.New(InvalidInput)
