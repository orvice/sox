package storage

import (
	"github.com/heatxsink/go-logstash"
)

type LogstashStorage struct {
	cli *logstash.Logstash
}

func NewLogStashStorage(addr string, port int) (*LogstashStorage, error) {
	l := logstash.New(addr, port, 30)
	_, err := l.Connect()
	if err != nil {
		return nil, err
	}
	return &LogstashStorage{
		cli: l,
	}, nil
}

func (l *LogstashStorage) Write(b []byte) error {
	return l.cli.Writeln(string(b))
}
