package tgff

import (
	"testing"
)

func TestParserReceiveOneSuccess(t *testing.T) {
	stream, abort := make(chan *token), make(chan bool)
	parser, _, _ := newParser(stream, abort)

	go func() {
		stream <- &token{numberToken, ""}
	}()

	_, err := parser.receiveOne(numberToken)

	assertSuccess(err, t)
}

func TestParserReceiveOneFailure(t *testing.T) {
	stream, abort := make(chan *token), make(chan bool)
	parser, _, _ := newParser(stream, abort)

	go func() {
		stream <- &token{identToken, ""}
	}()

	_, err := parser.receiveOne(numberToken)

	assertFailure(err, t)
}

func TestParserUnreceive(t *testing.T) {
	stream, abort := make(chan *token), make(chan bool)
	parser, _, _ := newParser(stream, abort)

	go func() {
		stream <- &token{numberToken, "First"}
		stream <- &token{numberToken, "Second"}
	}()

	token, _ := parser.receiveOne(numberToken)
	assertEqual(token.value, "First", t)

	parser.unreceive(token)

	token, _ = parser.receiveOne(numberToken)
	assertEqual(token.value, "First", t)

	token, _ = parser.receiveOne(numberToken)
	assertEqual(token.value, "Second", t)
}

func TestParserPeekOneOf(t *testing.T) {
	stream, abort := make(chan *token), make(chan bool)
	parser, _, _ := newParser(stream, abort)

	go func() {
		stream <- &token{numberToken, "First"}
	}()

	token, _ := parser.peekOneOf(numberToken, identToken)
	assertEqual(token.value, "First", t)

	token, _ = parser.receiveOne(numberToken)
	assertEqual(token.value, "First", t)
}

func TestParserRunClose(t *testing.T) {
	stream, abort := make(chan *token), make(chan bool)
	parser, success, failure := newParser(stream, abort)

	go parser.run()
	close(stream)

	select {
	case <-success:
	case <-failure:
		t.Fatal("expected a success")
	}
}
