package helper_test

import (
	exceptionsss_test "after-sales/api/expectionsss"
	"net/http"
)

func ReturnError(writer http.ResponseWriter, request *http.Request, err *exceptionsss_test.BaseErrorResponse) {
	if err.StatusCode == http.StatusUnauthorized {
		exceptionsss_test.NewAuthorizationException(writer, request, err)
		return
	} else if err.StatusCode == http.StatusBadRequest {
		exceptionsss_test.NewBadRequestException(writer, request, err)
		return
	} else if err.StatusCode == http.StatusUnprocessableEntity {
		exceptionsss_test.NewEntityException(writer, request, err)
		return
	} else if err.StatusCode == http.StatusNotFound {
		exceptionsss_test.NewNotFoundException(writer, request, err)
		return
	} else if err.StatusCode == http.StatusForbidden {
		exceptionsss_test.NewRoleException(writer, request, err)
		return
	} else if err.StatusCode == http.StatusConflict {
		exceptionsss_test.NewConflictException(writer, request, err)
		return
	} else {
		exceptionsss_test.NewAppException(writer, request, err)
		return
	}
}
