package parser

import "reflect"

type Reducer func(ctx *ReducerContext) interface{}

func ReduceAsString(ctx *ReducerContext) interface{} {
	ok, v := ctx.ListAsString()
	if !ok {
		panic("Cannot reduce as string")
	}
	return v
}

func ReReduce(ctx *ReducerContext) interface{} {
	return ctx.Reduce(ctx.Value.(Atom))
}

type ReducerContext struct {
	reducers *map[string]Reducer
	Value    interface{}
}

func (r ReducerContext) ListAsList() *AtomList {
	if v, ok := r.Value.(AtomList); ok {
		return &v
	}
	return nil
}

func (r ReducerContext) AtomList() []Atom {
	if v, ok := r.Value.(AtomList); ok {
		return v.value
	}
	return nil
}

func (r ReducerContext) ListAsString() (bool, string) {
	if v, ok := r.Value.(AtomList); ok {
		return ok, v.ReduceAsString()
	}
	return false, ""
}

func (r ReducerContext) Reduce(a Atom) interface{} {
	return ReduceInto(a, *r.reducers)
}

func (r ReducerContext) Flatten(v interface{}) interface{} {
	if reflect.TypeOf(v).Kind() != reflect.Slice {
		return v
	}
	var result []interface{}
	rv := reflect.ValueOf(v)
	for i := 0; i < rv.Len(); i++ {
		obj := rv.Index(i)
		switch obj.Type().Kind() {
		case reflect.Slice:
			result = append(result, r.Flatten(obj.Interface()).([]interface{})...)
		case reflect.Interface:
			innerObj := obj.Elem()
			if innerObj.Type().Kind() == reflect.Slice {
				result = append(result, r.Flatten(innerObj.Interface()).([]interface{})...)
			} else {
				result = append(result, obj.Interface())
			}
		default:
			result = append(result, obj.Interface())
		}
	}
	return result
}

func (r ReducerContext) IsNil(v interface{}) bool {
	k := reflect.TypeOf(v).Kind()
	if k == reflect.Ptr || k == reflect.Slice || k == reflect.Interface {
		return reflect.ValueOf(v).IsNil()
	}
	return v == nil
}

func (r ReducerContext) Iterate(i interface{}, fn func(i interface{})) {
	if r.IsNil(i) {
		return
	}

	if reflect.TypeOf(i).Kind() == reflect.Slice {
		v := reflect.ValueOf(i)
		for i := 0; i < v.Len(); i++ {
			obj := v.Index(i).Interface()
			if !r.IsNil(obj) {
				fn(obj)
			}
		}
	} else {
		fn(i)
	}
}

func (r ReducerContext) FindWithin(name string) Atom {
	if v := r.AtomList(); v != nil {
		for _, val := range v {
			if ref, ok := val.(RefResult); ok && ref.Name == name {
				return ref
			}
		}
	}
	return nil
}

func ReduceInto(root Atom, reducers map[string]Reducer) interface{} {
	switch v := root.(type) {
	case RefResult:
		if fn, ok := reducers[v.Name]; ok {
			return fn(&ReducerContext{
				reducers: &reducers,
				Value:    v.value,
			})
		}
		return v
	case AtomList:
		var result []interface{}
		for _, i := range v.value {
			result = append(result, ReduceInto(i, reducers))
		}
		return result
	}
	return root
}
