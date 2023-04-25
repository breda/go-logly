package http

import (
	"time"

	"github.com/breda/logly/internal/logly"
)

type HttpServer struct {
	logly *logly.Logly
}

type HttpError struct {
	Error string    `json:"error"`
	Time  time.Time `json:"time"`
}

type HttpAppendRequest struct {
	Data string `json:"data"`
}

type HttpAppendResponse struct {
	ID uint64 `json:"id"`
}

type HttpFetchRequest struct {
	ID uint64 `json:"id"`
}

type HttpFetchResponse struct {
	Data string `json:"data"`
}

type HttpIndexResponse struct {
	Message string `json:"message"`
}
