package store

import (
	"fmt"
	"testing"

	logv1 "github.com/breda/logly/api/v1"
)

func TestRead(t *testing.T) {
	store := NewInMemoryStore()

	store.records = append(store.records,
		&logv1.Record{Data: "log1"}, &logv1.Record{Data: "log2"}, &logv1.Record{Data: "log3"})

	secondRecord, err := store.Read(2)
	fmt.Println("Length of records: ", len(store.records))

	if err != nil {
		t.Fatal(err)
	}

	if secondRecord.Data != "log2" {
		t.FailNow()
	}
}

func TestWrite(t *testing.T) {
	store := NewInMemoryStore()

	written, err := store.Write(&logv1.Record{Data: "log1"})

	if err != nil {
		t.Fatal(err)
	}

	if written == 0 {
		t.FailNow()
	}

	if store.records[0].Data != "log1" {
		t.FailNow()
	}
}

func TestClear(t *testing.T) {
	store := NewInMemoryStore()

	store.Write(&logv1.Record{Data: "log1"})
	store.Write(&logv1.Record{Data: "log2"})
	store.Write(&logv1.Record{Data: "log3"})

	err := store.Clear()
	if err != nil {
		t.Fatal(err)
	}

	if len(store.records) > 0 {
		t.FailNow()
	}
}
