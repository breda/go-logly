package logly

import (
	"log"
	"sync"

	"github.com/breda/logly/internal/store"
)

type Logly struct {
	mtx   sync.Mutex
	store store.Store
}

func InMemory() *Logly {
	return &Logly{
		store: store.NewInMemoryStore(),
	}
}

func File() *Logly {
	fileStore, err := store.NewFileStore()
	if err != nil {
		log.Fatal(err)
	}

	return &Logly{
		store: fileStore,
	}
}
