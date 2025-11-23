package errors

type ErrorResponse struct {
	Error error `json:"error"`
}

func NewErrorResponse(e *Error) *ErrorResponse {
	return &ErrorResponse{
		Error: e,
	}
}
