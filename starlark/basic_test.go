package starlark

import (
	"testing"

	"github.com/drone/drone-go/plugin/config"
)

func TestExecNamedFunc(t *testing.T) {

	out, err := ExecNamedFunc("testdata/simple_test.star:linux_amd64", config.Request{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	buf, err := ValToYaml(out)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s\n", string(buf.Bytes()))
}

func TestReflectOnStruct(t *testing.T) {
	val := config.Request{}
	val.Build.ID = 234
	out := buildStarlarkVal(val.Build)
	t.Logf(out.String())
}
