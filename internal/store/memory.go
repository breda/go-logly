package store

import (
	"fmt"
	"sync"

	logv1 "github.com/breda/logly/api/v1"
	"github.com/breda/logly/internal/logger"
	"github.com/rs/zerolog"
)

type InMemoryStore struct {
	mtx    sync.Mutex
	Logger *zerolog.Logger

	size    uint64
	records []*logv1.Record
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		size:    0,
		records: make([]*logv1.Record, 0),
		Logger:  logger.New("memory-store"),
	}
}

func (s *InMemoryStore) Read(id uint64) (record *logv1.Record, err error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if id > s.size {
		return nil, fmt.Errorf("incorrect id %d: out of bounds", id)
	}

	record = s.records[id-1]
	err = nil
	return
}

func (s *InMemoryStore) Write(record *logv1.Record) (id uint64, err error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.records = append(s.records, record)
	s.size++

	id = s.size
	err = nil

	return
}

func (s *InMemoryStore) Clear() error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.records = make([]*logv1.Record, 0)
	s.size = 0

	return nil
}
