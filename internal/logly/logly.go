package logly

import (
	logv1 "github.com/breda/logly/api/v1"
)

func (l *Logly) Append(data string) (pos uint64, err error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	return l.store.Write(&logv1.Record{Data: data})
}

func (l *Logly) Fetch(pos uint64) (*logv1.Record, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	record, err := l.store.Read(pos)
	if err != nil {
		return nil, err
	}

	return record, nil
}
