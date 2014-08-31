package tgff

type tokenKind uint

type token struct {
	kind  tokenKind
	value string
}

const (
	errorToken tokenKind = iota
	blockCloseToken
	blockOpenToken
	controlToken
	identToken
	nameToken
	numberToken
	titleToken
)

func (k tokenKind) String() string {
	switch k {
	case errorToken:
		return "Error"
	case blockCloseToken:
		return "Block Close"
	case blockOpenToken:
		return "Block Open"
	case controlToken:
		return "Control"
	case identToken:
		return "Ident"
	case nameToken:
		return "Name"
	case numberToken:
		return "Number"
	case titleToken:
		return "Title"
	default:
		return "Unknown"
	}
}
