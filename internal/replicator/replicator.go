package replicator

import (
	"github.com/breda/logly/internal/logger"
	"github.com/rs/zerolog"
)

type ReplicationHandler struct {
	logger *zerolog.Logger
}

func NewReplicationHandler() *ReplicationHandler {
	return &ReplicationHandler{
		logger: logger.New("replicator"),
	}
}

func (h *ReplicationHandler) Join(name, addr string) error {
	return nil
}

func (h *ReplicationHandler) Leave(name string) error {
	return nil
}
