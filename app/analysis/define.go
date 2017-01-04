package analysis

type Type int

const (
	Normal Type = iota
	Func
)

var MethodMap = map[string]bool{
	"POST":   true,
	"GET":    true,
	"PUT":    true,
	"DELETE": true,
	"PATCH":  true,
}

const (
	Url    = "URL"
	Method = "METHOD"
	Header = "HEADER"
	Params = "PARAMS"
	Body   = "BODY"
	Count  = "COUNT"
)

var KeyMap = map[string]bool{
	Url:    true,
	Method: true,
	Header: true,
	Params: true,
	Body:   true,
	Count:  true,
}

const (
	SplitSymbol = "=>"
)
