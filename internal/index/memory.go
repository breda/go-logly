package index

import (
	"github.com/breda/logly/internal/logger"
	"github.com/rs/zerolog"
)

type InMemoryIndex struct {
	// Index maps the record ID to an offset in stored db file.
	index  map[int64]int64
	size   int64
	logger *zerolog.Logger
}

func InMemory() *InMemoryIndex {
	return &InMemoryIndex{
		index:  make(map[int64]int64),
		logger: logger.New("mem-index"),
	}
}

func (i *InMemoryIndex) Has(id int64) bool {
	_, found := i.index[id]
	return found
}

func (i *InMemoryIndex) Get(id int64) int64 {
	if !i.Has(id) {
		return 0
	}

	return i.index[id]
}

func (i *InMemoryIndex) Put(id, offset int64) {
	i.index[id] = offset
	i.size++
}

func (i *InMemoryIndex) Size() int64 {
	return i.size
}

func (i *InMemoryIndex) Type() string {
	return "inmem"
}

func (i *InMemoryIndex) Logger() *zerolog.Logger {
	return i.logger
}
