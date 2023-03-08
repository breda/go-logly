package store

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	logv1 "github.com/breda/logly/api/v1"
	"github.com/breda/logly/internal/index"
	"github.com/golang/protobuf/proto"
)

type FileStore struct {
	mtx sync.Mutex

	index index.Index

	nextId int64
	file   *os.File
}

const (
	RECORD_SIZE_WITDH_BYTES = 8
	DATA_FILE               = "data.db"
)

var (
	Encoding binary.ByteOrder = binary.BigEndian
)

func NewFileStore() (*FileStore, error) {
	file, err := os.OpenFile(DATA_FILE, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	// TODO: add a switch to use multiple index backends.
	indexSystem := index.NewBinaryTreeIndex()
	// indexSystem := index.InMemory()

	store := &FileStore{
		index:  indexSystem,
		file:   file,
		nextId: 1,
	}

	err = store.Init()
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (s *FileStore) Init() error {
	stat, err := s.file.Stat()
	if err != nil {
		return err
	}

	// If the file is new we do nothing.
	if stat.Size() == 0 {
		return nil
	}

	// Calculate the nextId
	var offset int64 = 0
	for {
		// Add the entry to the index
		s.index.Put(s.nextId, offset)

		sizeBytes := make([]byte, RECORD_SIZE_WITDH_BYTES)
		n, err := s.file.ReadAt(sizeBytes, offset)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		offset += int64(n)

		dataSize := Encoding.Uint64(sizeBytes)
		data := make([]byte, dataSize)

		n, err = s.file.ReadAt(data, offset)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		offset += int64(n)
		s.nextId++
	}

	log.Println("initialized db: next id:", s.nextId)
	log.Printf("initialized %s index with: %d entries", s.index.Type(), s.index.Size())

	return nil
}

func (s *FileStore) Write(record *logv1.Record) (id uint64, err error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	// Set the record ID from the one we got.
	record.Id = int64(s.nextId)

	// Then encode our record to binary
	data, err := proto.Marshal(record)
	if err != nil {
		return 0, err
	}

	// Encode the size of the record to binary
	sizeBytes := make([]byte, RECORD_SIZE_WITDH_BYTES)
	Encoding.PutUint64(sizeBytes, uint64(len(data)))

	// Write the size
	_, err = s.file.Write(sizeBytes)
	if err != nil {
		log.Fatal(err)
	}

	// Write the data
	_, err = s.file.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	// Advance the next ID
	s.nextId++

	return uint64(record.GetId()), nil
}

func (s *FileStore) Read(id uint64) (*logv1.Record, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if id > uint64(s.nextId) || id == 0 {
		return nil, fmt.Errorf("incorrect id %d", id)
	}

	if err := s.file.Sync(); err != nil {
		return nil, err
	}

	// Position ourselves at start of file.
	s.file.Seek(0, 0)

	var offset int64 = 0

	// Hey use the index!
	if s.index.Has(int64(id)) {
		offset = s.index.Get(int64(id))

		sizeBytes := make([]byte, RECORD_SIZE_WITDH_BYTES)
		s.file.ReadAt(sizeBytes, offset)

		dataSize := Encoding.Uint64(sizeBytes)
		data := make([]byte, dataSize)

		offset += RECORD_SIZE_WITDH_BYTES
		s.file.ReadAt(data, offset)

		var record logv1.Record
		err := proto.Unmarshal(data, &record)
		if err != nil {
			return nil, err
		}

		return &record, nil
	}

	for {
		sizeBytes := make([]byte, RECORD_SIZE_WITDH_BYTES)
		n, err := s.file.ReadAt(sizeBytes, offset)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		offset += int64(n)

		dataSize := Encoding.Uint64(sizeBytes)
		data := make([]byte, dataSize)

		n, err = s.file.ReadAt(data, offset)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		offset += int64(n)

		var record logv1.Record
		err = proto.Unmarshal(data, &record)
		if err != nil {
			return nil, err
		}

		if record.GetId() == int64(id) {
			return &record, nil
		}
	}

	return nil, fmt.Errorf("%d not found", id)
}

func (s *FileStore) Clear() error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.file.Truncate(0)
}
