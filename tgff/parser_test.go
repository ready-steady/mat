package tgff

import (
	"testing"
)

func TestParserRunEmpty(t *testing.T) {
	stream := make(chan *token)
	done := make(chan bool)

	parser, success, failure := newParser(stream, done)

	go parser.run()

	done <- true

	select {
	case <-success:
	case <-failure:
		t.Fatal("expected a success")
	}
}
