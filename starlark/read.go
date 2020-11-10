package starlark

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/drone/drone-go/plugin/config"
	"github.com/go-git/go-git/v5"
	"go.starlark.net/starlark"
)

func splitFuncPath(funcpath string) (string, string, error) {
	s := strings.SplitN(funcpath, ":", 2)
	if len(s) != 2 {
		return "", "", fmt.Errorf("Ill-formed filename:function string: %s", funcpath)
	}
	filename := s[0]
	funcname := s[1]

	return filename, funcname, nil
}

func ExecNamedFunc(funcpath string, wt *git.Worktree, req config.Request, extras map[string]string) (starlark.Value, error) {
	filename, funcname, err := splitFuncPath(funcpath)
	if err != nil {
		return err
	}

	thread := &starlark.Thread{
		Name: "rotocopter",
		Print: func(_ *starlark.Thread, msg string) {
			logrus.Info(msg)
		},
		Load: makeLoad(wt),
	}

	globals, err := parse(filename, wt, thread)
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
