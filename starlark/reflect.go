package starlark

import (
	"reflect"
	"strings"

	"github.com/drone/drone-go/plugin/config"
	"go.starlark.net/starlark"
)

func convertRequest(req config.Request, extra map[string]string) starlark.Value {
	out := starlark.NewDict(3)
	out.SetKey(starlark.String("build"), buildStarlarkVal(req.Build))
	out.SetKey(starlark.String("repo"), buildStarlarkVal(req.Repo))
	out.SetKey(starlark.String("config"), buildStringDict(extra))
	return out
}

func buildStarlarkVal(v interface{}) starlark.Value {
	e := reflect.ValueOf(v)
	out := starlark.NewDict(e.NumField())
	for i := 0; i < e.NumField(); i++ {
		fieldspec := e.Type().Field(i)
		tag, ok := fieldspec.Tag.Lookup("json")
		if !ok {
			tag = fieldspec.Name
		}
		if off := strings.Index(tag, ","); off >= 0 {
			tag = tag[:off]
		}
		var starlarkval starlark.Value
		switch fieldspec.Type.Kind() {
		case reflect.Bool:
			starlarkval = starlark.Bool(e.Field(i).Bool())
		case reflect.String:
			starlarkval = starlark.String(e.Field(i).String())
		case reflect.Int64:
			starlarkval = starlark.MakeInt64(int64(e.Field(i).Int()))
		case reflect.Int:
			starlarkval = starlark.MakeInt(int(e.Field(i).Int()))
		case reflect.Struct:
			fallthrough
		case reflect.Map:
			if mapval, ok := e.Field(i).Interface().(map[string]string); ok {
				starlarkval = buildStringDict(mapval)
			}
		default:
			continue
		}
		out.SetKey(
			starlark.String(tag),
			starlarkval,
		)
	}
	return out
}

func buildStringDict(from map[string]string) starlark.Value {
	out := starlark.NewDict(len(from))
	for k, v := range from {
		out.SetKey(starlark.String(k), starlark.String(v))
	}
	return out
}
