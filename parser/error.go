package parser

import (
	"fmt"
	"sort"
)

type ParseErrors []ParseError

func (p ParseErrors) Len() int               { return len(p) }
func (p ParseErrors) Less(i int, j int) bool { return p[i].Position > p[j].Position }
func (p ParseErrors) Swap(i int, j int)      { p[i], p[j] = p[j], p[i] }

type ParseError struct {
	Message  string
	Position int
	Errors   ParseErrors
}

func (p ParseError) Error() string {
	return fmt.Sprintf("%s at position %d", p.Message, p.Position)
}

func (p ParseError) Furthest() *ParseError {
	if len(p.Errors) == 0 {
		return &p
	}
	errs := p.Errors
	sort.Sort(errs)
	return &errs[0]
}

func (p *ParseError) Adopt(o ParseError) {
	p.Errors = append(p.Errors, o)
}

func Error(cur *Cursor, format string, args ...interface{}) *ParseError {
	return &ParseError{
		Message:  fmt.Sprintf(format, args...),
		Position: cur.pos + 1,
		Errors:   nil,
	}
}
