package utils

import (
	"errors"
	"strconv"
	"time"
)

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
