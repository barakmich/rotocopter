package starlark

import (
	"fmt"
	"io/ioutil"

	"go.starlark.net/starlark"
)

func parse(filename string, thread *starlark.Thread) (starlark.StringDict, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	globals, err := starlark.ExecFile(thread, filename, data, nil)
	if err != nil {
		return nil, err
	}
	return globals, nil
}

func args(v ...starlark.Value) starlark.Tuple {
	return starlark.Tuple(v)
}

// https://github.com/google/starlark-go/blob/4eb76950c5f02ec5bcfd3ca898231a6543942fd9/repl/repl.go#L175
func makeLoad() func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	type entry struct {
		globals starlark.StringDict
		err     error
	}

	var cache = make(map[string]*entry)

	return func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
		e, ok := cache[module]
		if e == nil {
			if ok {
				// request for package whose loading is in progress
				return nil, fmt.Errorf("cycle in load graph")
			}

			// Add a placeholder to indicate "load in progress".
			cache[module] = nil

			// Load it.
			thread := &starlark.Thread{Name: "exec " + module, Load: thread.Load}
			globals, err := starlark.ExecFile(thread, module, nil, nil)
			e = &entry{globals, err}

			// Update the cache.
			cache[module] = e
		}
		return e.globals, e.err
	}
}
