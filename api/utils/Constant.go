package utils

// "time"

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

// Etc
var LikeString string = "%%%s%%"

//	func BoolPtr(b bool) *bool {
//		return &b
//	}
func IntPtr(i int) *int {
	return &i
}

// func TimePtr(t time.Time) *time.Time {
// 	return &t
// }

func StringPtr(str string) *string {
	return &str
}
