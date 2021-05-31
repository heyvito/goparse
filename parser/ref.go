package parser

import (
	"context"
)

func ConsumerByRef(ctx context.Context, name string) Consumer {
	val := ctx.Value(ruleMapKey)
	if val == nil {
		panic("Invalid context")
	}
	mp, ok := val.(map[string]Consumer)
	if !ok {
		panic("BUG: Invalid ruleMapKey")
	}
	fn, ok := mp[name]
	if !ok {
		return nil
	}
	return fn
}

type RefConsumer struct {
	name string
}

func (o RefConsumer) Name() string   { return o.name }
func (o RefConsumer) String() string { return o.name }
func (RefConsumer) Weight() int      { return 0 }
func (o RefConsumer) TryConsume(ctx context.Context, c *Cursor) (Atom, error) {
	con := ConsumerByRef(ctx, o.name)
	if con == nil {
		return nil, Error(c, "unknown rule %s", o.name)
	}
	cd := c.dup()
	res := RefResult{parent: GetParent(ctx), Name: o.name}
	if v, err := con.TryConsume(SetParent(&res, ctx), &cd); err == nil {
		c.Merge(cd)
		res.value = v
		return res, nil
	} else {
		return nil, err
	}
}
