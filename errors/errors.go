package errors

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var (
	CLIENT_ID_REQUIRED  = &Error{Code: "CLIENT_ID_REQUIRED", Message: "Client Id is Required"}
	DUPLICATE_CLIENT_ID = &Error{Code: "DUPLICATE_CLIENT_ID", Message: "Client Id is Already in Use"}
	CLIENT_ID_NOT_FOUND = &Error{Code: "CLIENT_ID_NOT_FOUND", Message: "Client Id Not Found in the Room"}
	MESSAGE_IS_EMPTY    = &Error{Code: "MESSAGE_IS_EMPTY", Message: "Message Cannot be Empty"}
	REQUEST_TIMED_OUT   = &Error{Code: "REQUEST_TIMED_OUT", Message: "Request Timed Out"}
)
