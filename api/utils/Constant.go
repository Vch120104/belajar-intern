package utils

import (
	"errors"
	"strconv"
	"time"
)

// WO and Service Status
var SrvStatDraft int = 1    // 0 Draft
var SrvStatStart int = 2    // 10 Start
var SrvStatPending int = 3  // 15 Pending
var SrvStatStop int = 4     // 20 Stop
var SrvStatTransfer int = 5 // 25 Transfer
var SrvStatQcPass int = 6   // 30 QC Pass
var SrvStatReOrder int = 7  // 35 Re-Order

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
var LineTypeFee int = 5                // 5 Fee
var LinetypeAccesories int = 6         // 6 Accesories
var LinetypeConsumableMaterial int = 7 // 7 Consumable Material
var LineTypeSublet int = 8             // 8 Sublet
var LinetypeSouvenir int = 9           // 9 Souvenir

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
