package starlark

import (
	"fmt"
	"strings"

	"github.com/drone/drone-go/plugin/config"
	"go.starlark.net/starlark"
)

func ExecNamedFunc(f *f, req config.Request, extras map[string]string) (starlark.Value, error) {
	s := strings.SplitN(funcpath, ":", 2)
	if len(s) != 2 {
		return nil, fmt.Errorf("Ill-formed filename:function string: %s", funcpath)
	}
	filename := s[0]
	funcname := s[1]

	thread := &starlark.Thread{
		Name: "rotocopter",
		Print: func(_ *starlark.Thread, msg string) {
			fmt.Println(msg)
		},
		Load: makeLoad(),
	}

	globals, err := parse(filename, thread)
	if err != nil {
		return nil, err
	}
	v, ok := globals[funcname]
	if !ok {
		return nil, fmt.Errorf("No function named %s in %s", funcname, filename)

	}

	f, ok = v.(*starlark.Function)
	if !ok {
		return nil, fmt.Errorf("%s in %s is not a Function", funcname, filename)
	}

	arg := args(
		convertRequest(req, extras),
	)

	return starlark.Call(thread, f, arg, nil)
}
