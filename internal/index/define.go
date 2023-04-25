package index

import "github.com/rs/zerolog"

type Index interface {
	Has(id int64) bool
	Get(id int64) (offset int64)
	Put(id, offset int64)
	Size() int64
	Type() string
	Logger() *zerolog.Logger
}
