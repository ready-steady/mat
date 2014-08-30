package tgff

type parser struct {
	stream  <-chan token
	success chan<- *Result
	failure chan<- error
	result *Result
}

func newParser(stream <-chan token) (*parser, <-chan *Result, <-chan error) {
	success := make(chan *Result)
	failure := make(chan error)

	parser := &parser{
		stream:  stream,
		success: success,
		failure: failure,
		result:  &Result{},
	}

	return parser, success, failure
}

func (p *parser) run() {
	for token := range p.stream {
		switch token.kind {
		case errorToken:
			p.failure <- token.more[0].(error)
			return
		}
	}

	p.success <- p.result
}
