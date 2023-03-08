package http

import (
	"encoding/json"
	"net/http"

	"github.com/breda/logly/internal/logly"
)

func New(logly *logly.Logly) *HttpServer {
	return &HttpServer{
		logly: logly,
	}
}

func (s *HttpServer) HandleAppend(w http.ResponseWriter, r *http.Request) {
	var request HttpAppendRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	handleHttpError(err, w)

	pos, err := s.logly.Append(request.Data)
	handleHttpError(err, w)

	response := HttpAppendResponse{ID: pos}
	err = json.NewEncoder(w).Encode(&response)
	handleHttpError(err, w)
}

func (s *HttpServer) HandleFetch(w http.ResponseWriter, r *http.Request) {
	var request HttpFetchRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	handleHttpError(err, w)

	record, err := s.logly.Fetch(request.ID)
	handleHttpError(err, w)

	response := HttpFetchResponse{Data: record.Data}
	err = json.NewEncoder(w).Encode(&response)
	handleHttpError(err, w)
}
