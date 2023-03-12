package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/breda/logly/internal/logly"
)

func New(logly *logly.Logly) *HttpServer {
	server := &HttpServer{
		logly: logly,
	}

	return server
}

func (s *HttpServer) HandleAppend(w http.ResponseWriter, r *http.Request) {
	var request HttpAppendRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	handleHttpError(err, w, s.logly.Logger)

	pos, err := s.logly.Append(request.Data)
	handleHttpError(err, w, s.logly.Logger)

	response := HttpAppendResponse{ID: pos}
	err = json.NewEncoder(w).Encode(&response)
	handleHttpError(err, w, s.logly.Logger)
}

func (s *HttpServer) HandleFetch(w http.ResponseWriter, r *http.Request) {
	var request HttpFetchRequest

	// Get the ID from the URL query or from the request body
	if r.URL.Query().Has("id") {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		handleHttpError(err, w, s.logly.Logger)

		request.ID = uint64(id)
	} else {
		err := json.NewDecoder(r.Body).Decode(&request)
		handleHttpError(err, w, s.logly.Logger)
	}

	record, err := s.logly.Fetch(request.ID)
	handleHttpError(err, w, s.logly.Logger)

	if record != nil {
		response := HttpFetchResponse{Data: record.Data}
		err = json.NewEncoder(w).Encode(&response)
		handleHttpError(err, w, s.logly.Logger)
	}
}

func (s *HttpServer) HandleIndex(w http.ResponseWriter, r *http.Request) {
	response := HttpIndexResponse{Message: "Hello, this is Logly"}
	err := json.NewEncoder(w).Encode(&response)
	handleHttpError(err, w, s.logly.Logger)
}
