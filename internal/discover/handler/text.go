package handler

import (
	"fmt"

	"github.com/breda/logly/internal/logger"
	"github.com/rs/zerolog"
)

type LogDiscoveryHandler struct {
	logger *zerolog.Logger
}

func NewLogHandler() *LogDiscoveryHandler {
	return &LogDiscoveryHandler{
		logger: logger.New("discover-log"),
	}
}

func (h *LogDiscoveryHandler) Join(name, addr string) error {
	h.logger.Info().Msg(fmt.Sprintf("member %s with address %s joined the cluster", name, addr))
	return nil
}

func (h *LogDiscoveryHandler) Leave(name string) error {
	h.logger.Info().Msg(fmt.Sprintf("member %s left the cluster", name))
	return nil
}
