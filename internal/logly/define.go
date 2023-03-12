package logly

import (
	"log"
	"sync"

	"github.com/breda/logly/internal/logger"
	"github.com/breda/logly/internal/store"
	"github.com/rs/zerolog"
)

type Logly struct {
	mtx    sync.Mutex
	Logger *zerolog.Logger
	store  store.Store
}

func InMemory() *Logly {
	return &Logly{
		store:  store.NewInMemoryStore(),
		Logger: logger.New("logly"),
	}
}

func File() *Logly {
	fileStore, err := store.NewFileStore()
	if err != nil {
		log.Fatal(err)
	}

	return &Logly{
		store:  fileStore,
		Logger: logger.New("logly"),
	}
}
