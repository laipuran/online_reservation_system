package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JSON(w http.ResponseWriter, httpStatus int, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(resp)
}

func OK(data interface{}) Response {
	return Response{Code: 200, Message: "ok", Data: data}
}

func Created(data interface{}) Response {
	return Response{Code: 201, Message: "created", Data: data}
}

func Fail(msg string) Response {
	return Response{Code: 400, Message: msg, Data: nil}
}

func Unauthorized(msg string) Response {
	return Response{Code: 401, Message: msg, Data: nil}
}

func Conflict(msg string) Response {
	return Response{Code: 409, Message: msg, Data: nil}
}

func ServerError(msg string) Response {
	return Response{Code: 500, Message: msg, Data: nil}
}
