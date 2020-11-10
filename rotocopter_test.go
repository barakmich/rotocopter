package main

import (
	"context"
	"strings"
	"testing"

	"github.com/barakmich/rotocopter/roto"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/config"
)

func TestLoadGit(t *testing.T) {
	conf := roto.Config{
		GitYamlRepo: "testdata/",
	}
	req := config.Request{
		Repo: drone.Repo{
			Slug: "barak/testpipeline",
		},
	}

	r := roto.New(conf)
	config, err := r.Find(context.Background(), &req)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(config.Data), "Whatup") {
		t.Fatal("data doesn't contain echo command")
	}
	t.Log(config.Data)
}
