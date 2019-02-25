package api

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

//Response root element of the http responses.
type Response struct {
	Data  *interface{} `json:"data,omitempty"`
	Error *Error       `json:"error,omitempty"`
	Links *Links       `json:"links,omitempty"`
}

//Error contains the error information
type Error struct {
	InternalCode int    `json:"code"`
	Message      string `json:"msg"`
}

//Links contains
type Links struct {
	Self string `json:"self"`
}

//SendResponse converts a code and a interface into a Response and sends the response.
func SendResponse(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	var rspPayload = &Response{
		Data:  &payload,
		Error: nil,
		Links: &Links{Self: fmt.Sprintf("%s%s", r.Host, r.URL.String())},
	}
	sendJSONResponse(w, r, code, rspPayload)
}

//SendErrorResponse converts a code and an error into a Response with Error not nil and sends the response.
func SendErrorResponse(w http.ResponseWriter, r *http.Request, code int, err error) {
	var rspPayload = &Response{
		Data:  nil,
		Error: &Error{InternalCode: code, Message: err.Error()},
		Links: &Links{Self: fmt.Sprintf("%s%s", r.Host, r.URL.String())},
	}
	logrus.WithError(err).Info("failed response")
	sendJSONResponse(w, r, code, rspPayload)
}

func sendJSONResponse(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(payload)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(code)
	_, _ = w.Write(response)
}
