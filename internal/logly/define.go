package logly

import (
	"log"
	"net"
	"sync"

	"github.com/breda/logly/internal/discover"
	"github.com/breda/logly/internal/discover/handler"
	"github.com/breda/logly/internal/logger"
	"github.com/breda/logly/internal/store"
	"github.com/rs/zerolog"
)

type Logly struct {
	mtx        sync.Mutex
	Logger     *zerolog.Logger
	store      store.Store
	membership discover.Membership
}

func getMembership(name, addr, port, cluster string) (*discover.Membership, error) {
	tags := make(map[string]string)
	tags["rpc_addr"] = net.JoinHostPort(addr, port)

	config := discover.Config{
		NodeName: name,
		BindAddr: net.JoinHostPort(addr, port),
		Tags:     tags,
	}

	if cluster != "" {
		config.StartJoinAddrs = make([]string, 1)
		config.StartJoinAddrs[0] = cluster
	}

	return discover.New(handler.NewLogHandler(), config)
}

func InMemory(name, addr, port, cluster string) *Logly {
	membership, err := getMembership(name, addr, port, cluster)
	if err != nil {
		panic(err)
	}

	return &Logly{
		store:      store.NewInMemoryStore(),
		Logger:     logger.New("logly"),
		membership: *membership,
	}
}

func File(name, addr, port, cluster string) *Logly {
	membership, err := getMembership(name, addr, port, cluster)
	if err != nil {
		panic(err)
	}

	fileStore, err := store.NewFileStore()
	if err != nil {
		log.Fatal(err)
	}

	return &Logly{
		store:      fileStore,
		Logger:     logger.New("logly"),
		membership: *membership,
	}
}
