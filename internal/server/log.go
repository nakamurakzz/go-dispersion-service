package server

import (
	"fmt"
	"log"
	"sync"
)

type Log struct {
	mu      sync.Mutex
	records []Record
}

type Record struct {
	Value  []byte `json:"value"`
	Offset int64  `json:"offset"`
}

func NewLog() *Log {
	return &Log{}
}

// スライスにレコードを追加する
func (l *Log) Append(record Record) (int64, error) {
	log.Println("Append: ", record)
	l.mu.Lock()
	defer l.mu.Unlock()

	record.Offset = int64(len(l.records))
	l.records = append(l.records, record)
	return record.Offset, nil
}

// インデックスを指定してスライスからレコードを取得する
func (l *Log) Read(offset int64) (Record, error) {
	log.Println("Read: ", offset)
	l.mu.Lock()
	defer l.mu.Unlock()

	if int(offset) >= len(l.records) {
		return Record{}, ErrOffsetNotFound
	}
	return l.records[offset], nil
}

var ErrOffsetNotFound = fmt.Errorf("offset not found")
