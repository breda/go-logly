package store

import (
	logv1 "github.com/breda/logly/api/v1"
)

type Store interface {
	Read(id uint64) (record *logv1.Record, err error)
	Write(record *logv1.Record) (id uint64, err error)
	Clear() error
}
