package ingest


import (
	"fmt"
	"strings"
)


/*
The grammar

Format -> Exprs ;

Exprs -> Expr Exprs
       | Expr
       ;

Expr -> Var
      | LITERAL
      ;

Var -> DOLLAR LPAREN Name RPAREN
     ;

Name -> LITERAL Name
      | LITERAL
      ;
*/

const (
	FormatChar = 1 << iota
	FormatVar
)

type Format []FormatElement

func (f Format) String() string {
	parts := make([]string, len(f))
	for _, e := range f {
		parts = append(parts, e.String())
	}
	return strings.Join(parts, "")
}

func (f Format) VerboseString() string {
	parts := make([]string, len(f))
	for _, e := range f {
		parts = append(parts, e.VerboseString())
	}
	return strings.Join(parts, "")
}

type FormatElement struct {
	Type uint
	Name string
	Char byte
}

func (fe FormatElement) String() string {
	switch fe.Type {
	case FormatChar:
		return string([]byte{fe.Char})
	case FormatVar:
		return fmt.Sprintf("$(%v)", fe.Name)
	default:
		panic(fmt.Errorf("unexpect format element, %v", fe.Type))
	}
}

func (fe FormatElement) VerboseString() string {
	switch fe.Type {
	case FormatChar:
		return fmt.Sprintf("<char %v>", string([]byte{fe.Char}))
	case FormatVar:
		return fmt.Sprintf("<var %v>", fe.Name)
	default:
		panic(fmt.Errorf("unexpect format element, %v", fe.Type))
	}
}

func (f Format) Validate() error {
	for i, e := range f {
		var err error
		switch e.Type {
		case FormatChar:
		case FormatVar:
			if i + 1 < len(f) && f[i+1].Type == FormatVar {
				err = fmt.Errorf("variables must be seperated by a constant, '%v' '%v'", e, f[i+1])
			}
		default:
			err = fmt.Errorf("unexpect format element, %v", e)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (f Format) Parse(bytes []byte) (map[string]string, error) {
	meta := make(map[string]string, len(f)/2 + 1)
	err := f.ParseInto(bytes, meta)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

func (f Format) ParseInto(bytes []byte, meta map[string]string) error {
	err := f.Validate()
	if err != nil {
		return err
	}
	j := 0
	for i, e := range f {
		var err error
		switch e.Type {
		case FormatChar:
			j, err = f.scan_char(i, j, bytes)
		case FormatVar:
			j, err = f.scan_var(i, j, bytes, meta)
		default:
			err = fmt.Errorf("unexpect format element, %v", e)
		}
		if err != nil {
			return err
		}
	}
	if j != len(bytes) {
		return fmt.Errorf("unconsumed input at end '%v'", string(bytes[j:]))
	}
	return nil
}

func (f Format) scan_char(i, j int, bytes []byte) (int, error) {
	if j >= len(bytes) {
		return j, fmt.Errorf("unexpected EOF, expected '%v'", f[i])
	}
	if bytes[j] != f[i].Char {
		return j, fmt.Errorf("unexpected '%v', expected '%v'", string(bytes[j]), string(f[i].Char))
	}
	return j+1, nil
}

func (f Format) scan_var(i, j int, bytes []byte, meta map[string]string) (int, error) {
	var eof bool
	var stop byte
	if i + 1 >= len(f) {
		eof = true
	} else {
		if f[i+1].Type != FormatChar {
			return j, fmt.Errorf("Format string invalid, to variables next to each other with out const separator")
		}
		stop = f[i+1].Char
	}
	buf := make([]byte, 0, len(bytes)-j)
	c := j
	for ; c < len(bytes); c++ {
		if !eof && bytes[c] == stop {
			break
		}
		buf = append(buf, bytes[c])
	}
	if len(buf) == 0 {
		return j, fmt.Errorf("Varaible %s not supplied", f[i].Name)
	}
	meta[f[i].Name] = string(buf)
	return c, nil
}

type Consumer interface {
	Consume(i int) (int, interface{}, error)
}

type StrProduction struct {
	name string
	Productions map[string]Consumer
}

func (p *StrProduction) Consume(i int) (int, interface{}, error) {
	return p.Productions[p.name].Consume(i)
}

type FnProduction func(i int) (int, interface{}, error)

func (fn FnProduction) Consume(i int) (int, interface{}, error) {
	return fn(i)
}

func ParseFormatString(format string) (Format, error) {
	P := make(map[string]Consumer)
	S := func(name string) Consumer {
		return &StrProduction{
			name: name,
			Productions: P,
		}
	}

	var (
		Consume func(byte) Consumer
		Concat func(...Consumer) (func(func(...interface{})(interface{},error)) Consumer)
		Alt func(...Consumer) Consumer
		LITERAL Consumer
	)

	Consume = func(token byte) Consumer {
		return FnProduction(func(i int) (int, interface{}, error) {
			if i == len(format) {
				return i, nil, fmt.Errorf("Ran off end of input. Expected '%v'", token)
			}
			if format[i] == token {
				return i+1, format[i], nil
			}
			return i, nil, fmt.Errorf("Expected %v got %v", string([]byte{token}), format[i:i+1])
		})
	}

	LITERAL = FnProduction(func(i int) (int, interface{}, error) {
		if i == len(format) {
			return i, nil, fmt.Errorf("Ran off end of input. Expected any char")
		}
		return i+1, FormatElement{Type:FormatChar, Char:format[i]}, nil
	})

	Concat = func(consumers ...Consumer) func(func(...interface{})(interface{}, error)) Consumer {
		return func(action func(...interface{})(interface{},error)) Consumer {
			return FnProduction(func(i int) (int, interface{}, error) {
				var nodes []interface{}
				var n interface{}
				var err error
				j := i
				for _, c := range consumers {
					j, n, err = c.Consume(j)
					if err == nil {
						nodes = append(nodes, n)
					} else {
						return i, nil, err
					}
				}
				an, aerr := action(nodes...)
				if aerr != nil {
					return i, nil, aerr
				}
				return j, an, nil
			})
		}
	}

	Alt = func(consumers ...Consumer) Consumer {
		return FnProduction(func(i int) (int, interface{}, error) {
			var err error
			for _, c := range consumers {
				j, n, e := c.Consume(i)
				if e == nil {
					return j, n, nil
				} else {
					err = e
				}
			}
			return i, nil, err
		})
	}

	P["Format"] = Concat(S("Exprs"))(func(nodes ...interface{}) (interface{}, error) {
		rf := nodes[0].(Format)
		f := make(Format, 0, len(rf))
		for i := len(rf) - 1; i >= 0; i-- {
			f = append(f, rf[i])
		}
		return f, nil
	})

	P["Exprs"] = Alt(
		Concat(S("Expr"), S("Exprs"))(func(nodes ...interface{}) (interface{}, error) {
			expr := nodes[0].(FormatElement)
			exprs := nodes[1].(Format)
			return append(exprs, expr), nil
		}),
		Concat(S("Expr"))(func(nodes ...interface{}) (interface{}, error) {
			expr := nodes[0].(FormatElement)
			elems := make(Format, 0, len(format))
			elems = append(elems, expr)
			return elems, nil
		}),
	)

	P["Expr"] = Alt(
		S("Var"),
		LITERAL,
	)

	P["Var"] = Concat(Consume('$'), Consume('('), S("Name"), Consume(')'))(
		func(nodes ...interface{}) (interface{}, error) {
			name := nodes[2].(string)
			fe := FormatElement{Type:FormatVar, Name:name}
			return fe, nil
		})

	P["Name"] = FnProduction(func(i int) (int, interface{}, error) {
		buf := make([]byte, 0, 10)
		for j := i; j < len(format); j++ {
			if format[j] == ')' {
				return j, string(buf), nil
			}
			buf = append(buf, format[j])
		}
		return i, nil, fmt.Errorf("ran off the end of the input")
	})

	i, node, err := P["Format"].Consume(0)
	if err != nil {
		return nil, err
	}
	if len(format) != i {
		return nil, fmt.Errorf("unconsumed input %v", format[i:])
	}
	return node.(Format), nil
}

