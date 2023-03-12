package logly

import (
	logv1 "github.com/breda/logly/api/v1"
)

func (l *Logly) Append(data string) (id uint64, err error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	id, err = l.store.Write(&logv1.Record{Data: data})
	if err != nil {
		l.Logger.Fatal().Err(err).Send()
	}

	l.Logger.Debug().Msg("successfully appended a new log entry")
	return
}

func (l *Logly) Fetch(id uint64) (*logv1.Record, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	record, err := l.store.Read(id)
	if err != nil {
		return nil, err
	}

	l.Logger.Debug().Str("message", "fetch record request completed").Uint64("id", id).Send()
	return record, nil
}
