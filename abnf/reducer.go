package abnf

import p "github.com/heyvito/goparse/parser"

var Reducer = map[string]p.Reducer{
	"rulelist": func(ctx *p.ReducerContext) interface{} {
		var result []Rule
		for _, i := range ctx.AtomList() {
			result = append(result, ctx.Reduce(i).(Rule))
		}
		return &RuleList{Rules: result}
	},
	"rule": func(ctx *p.ReducerContext) interface{} {
		return Rule{
			Name:      ctx.Reduce(ctx.FindWithin("rulename")).(RuleName),
			DefinedAs: ctx.Reduce(ctx.FindWithin("defined-as")).(DefinedAs),
			Elements:  ctx.Reduce(ctx.FindWithin("elements")).(Elements),
		}
	},
	"rulename": func(ctx *p.ReducerContext) interface{} {
		_, v := ctx.ListAsString()
		return RuleName{Name: v}
	},
	"defined-as": func(ctx *p.ReducerContext) interface{} {
		_, v := ctx.ListAsString()
		return DefinedAs{Value: v}
	},
	"elements": func(ctx *p.ReducerContext) interface{} {
		return Elements{Alternation: ctx.Reduce(ctx.Value.(p.AtomList).First()).(Alternation)}
	},
	"alternation": func(ctx *p.ReducerContext) interface{} {
		var opts []interface{}
		for _, v := range ctx.AtomList() {
			if red := ctx.Reduce(v); !ctx.IsNil(red) {
				opts = append(opts, red)
			}
		}
		var cats []Concatenation
		for _, i := range ctx.Flatten(opts).([]interface{}) {
			cats = append(cats, i.(Concatenation))
		}

		return Alternation{Elements: cats}
	},
	"concatenation": func(ctx *p.ReducerContext) interface{} {
		// Concatenation contains a single repetition, followed by a
		// list of repetitions. Here we may attempt to fix that.
		var opts []Repetition
		for _, v := range ctx.AtomList() {
			ctx.Iterate(ctx.Flatten(ctx.Reduce(v)), func(i interface{}) {
				opts = append(opts, i.(Repetition))
			})
		}

		return Concatenation{Elements: opts}
	},
	"repeat": func(ctx *p.ReducerContext) interface{} {
		list := ctx.ListAsList()
		var min, max int
		switch list.Len() {
		case 1:
			_, v := ctx.ListAsList().ReduceAsInt()
			min = v
			max = v
		case 3:
			_, min = list.Nth(0).(p.AtomList).ReduceAsInt()
			_, max = list.Nth(2).(p.AtomList).ReduceAsInt()
		}

		return Repeat{
			Min: min,
			Max: max,
		}
	},
	"repetition": func(ctx *p.ReducerContext) interface{} {
		list := ctx.ListAsList()
		if list.Len() != 2 {
			panic("BUG: Invalid repetition format")
		}

		v := ctx.Reduce(list.Nth(1)).(Element)
		rp := list.First().(p.OptionVal)
		if !rp.Valid {
			return Repetition{
				Meta:    nil,
				Element: v,
			}
		}
		meta := ctx.Reduce(rp.Value().(p.Atom)).(Repeat)
		return Repetition{
			Meta:    &meta,
			Element: v,
		}
	},
	"element": func(ctx *p.ReducerContext) interface{} {
		reduced := ctx.Reduce(ctx.Value.(p.Atom))
		return Element{Inner: reduced.(Node)}
	},
	"group": func(ctx *p.ReducerContext) interface{} {
		list := ctx.ListAsList()
		if list.Len() != 3 {
			panic("BUG: Invalid group format")
		}
		return Group{Elements: ctx.Reduce(list.Nth(1)).(Alternation)}
	},
	"option": func(ctx *p.ReducerContext) interface{} {
		list := ctx.ListAsList()
		if list.Len() != 3 {
			panic("BUG: Invalid option format")
		}
		return Option{
			Elements: ctx.Reduce(list.Nth(1)).(Alternation),
		}
	},
	"char-val": func(ctx *p.ReducerContext) interface{} {
		return CharVal{Value: ctx.ListAsList().Nth(1).(p.AtomList).ReduceAsString()}
	},
	"num-val": func(ctx *p.ReducerContext) interface{} {
		return ctx.Reduce(ctx.ListAsList().Nth(1))
	},
	"prose-val": p.ReduceAsString,
	"hex-val": func(ctx *p.ReducerContext) interface{} {
		opt := ctx.ListAsList().Nth(2).(p.OptionVal)
		_, single := ctx.ListAsList().Nth(1).(p.AtomList).ReduceAsHex()

		result := Numeric{
			Mode:   NumericModeSingle,
			Single: single,
		}
		if opt.Valid {
			optVal := opt.Value().(p.AtomList)
			if optVal.Nth(0).(p.Char).Value() == "-" {
				result.Mode = NumericModeRange
				_, to := optVal.Nth(1).(p.AtomList).ReduceAsHex()
				result.Range = Range{
					From: single,
					To:   to,
				}
			} else {
				result.Mode = NumericModeSequence
			}
		}

		return HexVal{result}
	},
	"dec-val":func(ctx *p.ReducerContext) interface{} {
		opt := ctx.ListAsList().Nth(2).(p.OptionVal)
		_, single := ctx.ListAsList().Nth(1).(p.AtomList).ReduceAsInt()

		result := Numeric{
			Mode:   NumericModeSingle,
			Single: single,
		}
		if opt.Valid {
			optVal := opt.Value().(p.AtomList)
			if optVal.Nth(0).(p.Char).Value() == "-" {
				result.Mode = NumericModeRange
				_, to := optVal.Nth(1).(p.AtomList).ReduceAsInt()
				result.Range = Range{
					From: single,
					To:   to,
				}
			} else {
				result.Mode = NumericModeSequence
			}
		}

		return DecVal{result}
	},
}
