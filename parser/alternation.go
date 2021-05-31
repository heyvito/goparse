package parser

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

type WeightedResult struct {
	c Cursor
	v Atom
	w int
}

type WeightedResults []WeightedResult

func (w WeightedResults) Len() int           { return len(w) }
func (w WeightedResults) Less(i, j int) bool { return w[i].w > w[j].w }
func (w WeightedResults) Swap(i, j int)      { w[i], w[j] = w[j], w[i] }

type AlternationConsumer struct {
	cons []Consumer
}

func (a AlternationConsumer) Name() string { return "ALTERNATION" }
func (AlternationConsumer) Weight() int    { return 1 }
func (a AlternationConsumer) String() string {
	str := make([]string, len(a.cons))
	for i, c := range a.cons {
		str[i] = c.String()
	}
	return fmt.Sprintf("( %s )", strings.Join(str, " / "))
}
func (a AlternationConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	var results WeightedResults
	var errors ParseErrors
	for _, v := range a.cons {
		cd := c.dup()
		if ret, err := v.TryConsume(ctx, &cd); err == nil {
			results = append(results, WeightedResult{cd, ret, v.Weight()})
		} else {
			errors = append(errors, *err.(*ParseError))
		}
	}

	if len(results) == 0 {
		err := Error(c, "Expected one of the following rules to be met, but none could be matched:")
		for _, e := range errors {
			err.Adopt(e)
		}

		return nil, err.Furthest()
	}

	sort.Sort(results)
	res := results[0]
	c.Merge(res.c)
	return res.v, nil
}
