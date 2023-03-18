package logly

import (
	"time"

	logv1 "github.com/breda/logly/api/v1"
)

func (l *Logly) Append(data string) (id uint64, err error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	start := time.Now()
	id, err = l.store.Write(&logv1.Record{Data: data})
	if err != nil {
		l.Logger.Fatal().Err(err).Send()
	}

	elapsed := time.Since(start)
	l.Logger.
		Debug().
		Str("message", "successfully appended a new log entry").
		Str("took", elapsed.String()).
		Send()
	return
}

func (l *Logly) Fetch(id uint64) (*logv1.Record, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	start := time.Now()
	record, err := l.store.Read(id)
	if err != nil {
		return nil, err
	}

	elapsed := time.Since(start)
	l.Logger.
		Debug().
		Str("message", "fetch record request completed").
		Uint64("id", id).
		Str("took", elapsed.String()).
		Send()

	return record, nil
}
