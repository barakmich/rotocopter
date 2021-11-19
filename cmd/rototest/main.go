package main

import (
	"fmt"
	"os"

	"github.com/barakmich/rotocopter/starlark"
)

func getYamlFrom(v string) (string, error) {
	starlarkval, err := starlark.ExecFuncFromFile(os.ExpandEnv("$PWD"), v, nil)
	if err != nil {
		return "", err
	}
	buf, err := starlark.ValToYaml(starlarkval)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func main() {
	s, err := getYamlFrom(os.Args[1])
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("%s", s)
}
