package starlark

import (
	"testing"

	"github.com/drone/drone-go/plugin/config"
)

func TestReflectOnStruct(t *testing.T) {
	val := config.Request{}
	val.Build.ID = 234
	out := buildStarlarkVal(val.Build)
	t.Logf(out.String())
}
