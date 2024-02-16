package types

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

type Response struct {
	Response string `json:"response"`
	Message  string `json:"message"`
}

func NewResponse(response string, message string) Response {
	return Response{
		Response: response,
		Message:  message,
	}

}
