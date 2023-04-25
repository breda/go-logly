package discover

import (
	"github.com/breda/logly/internal/logger"
	"github.com/hashicorp/serf/serf"
	"github.com/rs/zerolog"
)

type Config struct {
	NodeName       string
	BindAddr       string
	Tags           map[string]string
	StartJoinAddrs []string
}

type Handler interface {
	Join(name, addr string) error
	Leave(name string) error
}

type Membership struct {
	Config
	handler Handler
	serf    *serf.Serf
	events  chan serf.Event
	Logger  *zerolog.Logger
}

func New(handler Handler, config Config) (*Membership, error) {
	m := &Membership{
		Config:  config,
		handler: handler,
		Logger:  logger.New("discovery"),
	}

	if err := m.setupSerf(); err != nil {
		return nil, err
	}

	return m, nil
}
