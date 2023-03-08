package http

import (
	"github.com/breda/logly/internal/logly"
)

type HttpServer struct {
	logly *logly.Logly
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
