package parser

import (
	"context"
	"fmt"
	"strings"
)

type RepetitionMode int

const (
	RepeatPlus RepetitionMode = iota
	RepeatStar
	RepeatMin
	RepeatMinMax
)

type RepetitionConsumer struct {
	mode RepetitionMode
	min  int
	max  int
	con  Consumer
}

func (RepetitionConsumer) Weight() int { return 0 }

func (r RepetitionConsumer) String() string {
	str := strings.Builder{}
	switch r.mode {
	case RepeatPlus:
		str.WriteRune('+')
	case RepeatStar:
		str.WriteRune('*')
	case RepeatMin:
		str.WriteString(fmt.Sprintf("%d*", r.min))
	case RepeatMinMax:
		str.WriteString(fmt.Sprintf("%d*%d", r.min, r.max))
	}
	str.WriteString(r.con.String())
	return str.String()
}

func (r RepetitionConsumer) Name() string {
	str := strings.Builder{}
	str.WriteString("Repeat(")
	switch r.mode {
	case RepeatPlus:
		str.WriteRune('+')
	case RepeatStar:
		str.WriteRune('*')
	case RepeatMin:
		str.WriteString(fmt.Sprintf("%d*", r.min))
	case RepeatMinMax:
		str.WriteString(fmt.Sprintf("%d*%d", r.min, r.max))
	}
	str.WriteString(", ")
	str.WriteString(r.con.String())
	str.WriteRune(')')
	return str.String()
}

func (r RepetitionConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	result := AtomList{parent: GetParent(ctx)}
	var list []Atom
	cd := c.dup()
	switch r.mode {
	case RepeatPlus:
		for {
			v, err := r.con.TryConsume(SetParent(&result, ctx), &cd)
			if err == nil {
				if v != nil {
					list = append(list, v)
				}
				continue
			}
			if len(list) > 0 {
				c.Merge(cd)
			} else {
				return result, err
			}
			result.value = list
			return result, nil
		}
	case RepeatStar:
		moveCursor := false
		for {
			// Do we have CD?
			if ok, _ := cd.TryPeek(); ok {
				v, err := r.con.TryConsume(SetParent(&result, ctx), &cd)
				if err == nil {
					moveCursor = true
					if v != nil {
						list = append(list, v)
					}
					continue
				}
			}
			if moveCursor {
				c.Merge(cd)
			}
			result.value = list
			return result, nil
		}
	case RepeatMin:
		for {
			v, err := r.con.TryConsume(SetParent(&result, ctx), &cd)
			if err == nil {
				if v != nil {
					list = append(list, v)
				}
				continue
			}
			if len(list) == r.min {
				c.Merge(cd)
			} else {
				return nil, err
			}
			result.value = list
			return result, nil
		}
	case RepeatMinMax:
		for {
			v, err := r.con.TryConsume(SetParent(&result, ctx), &cd)
			if err == nil {
				if v != nil {
					list = append(list, v)
				}
				continue
			}
			if len(list) >= r.min && len(list) <= r.max {
				c.Merge(cd)
			} else {
				return nil, err
			}
			result.value = list
			return result, nil
		}
	default:
		panic("Invalid repetition mode")
	}
}
