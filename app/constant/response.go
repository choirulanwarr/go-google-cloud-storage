package constant

import "net/http"

const RequestIDKey = "API_CALL_ID"

type ResponseMap struct {
	Code    int
	Status  string
	Message string
}

const (
	ResponseStatusSuccess = "SUCCESS"
	ResponseStatusFailed  = "FAILED"
)

var (
	Res200Save = ResponseMap{Code: http.StatusOK, Status: ResponseStatusSuccess, Message: "Success save data"}
	Res200Get  = ResponseMap{Code: http.StatusOK, Status: ResponseStatusSuccess, Message: "Success get data"}
)

var (
	Res400InvalidPayload     = ResponseMap{Code: http.StatusBadRequest, Status: ResponseStatusFailed, Message: "Invalid payload data"}
	Res400FailedDataNotFound = ResponseMap{Code: http.StatusBadRequest, Status: ResponseStatusFailed, Message: "Data not found"}
)

var (
	Res422SomethingWentWrong = ResponseMap{Code: http.StatusUnprocessableEntity, Status: ResponseStatusFailed, Message: "Something Went Wrong"}
)
