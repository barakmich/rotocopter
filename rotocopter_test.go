package main

import (
	"context"
	"strings"
	"testing"

	"github.com/barakmich/rotocopter/roto"
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/config"
)

func loadAndRunTestdata(t *testing.T, req *config.Request) *drone.Config {
	conf := roto.Config{
		GitYamlRepo: "testdata/",
	}
	r := roto.New(conf)
	config, err := r.Find(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	return config
}

func TestLoadGit(t *testing.T) {
	req := config.Request{
		Repo: drone.Repo{
			Slug: "barak/testpipeline",
		},
	}
	config := loadAndRunTestdata(t, &req)
	if !strings.Contains(string(config.Data), "Whatup") {
		t.Fatal("data doesn't contain echo command")
	}
	t.Log(config.Data)
}

func TestBasicStarlark(t *testing.T) {
	req := config.Request{
		Repo: drone.Repo{
			Slug: "barak/teststarlark",
		},
	}

	config := loadAndRunTestdata(t, &req)
	if !strings.Contains(string(config.Data), "kind: docker") {
		t.Fatal("data doesn't contain docker kind")
	}
	t.Log(config.Data)
}

func TestStarlarkThatLoads(t *testing.T) {
	req := config.Request{
		Repo: drone.Repo{
			Slug: "barak/metatest",
		},
	}

	config := loadAndRunTestdata(t, &req)
	if !strings.Contains(string(config.Data), "kind: docker") {
		t.Fatal("data doesn't contain docker kind")
	}
	if !strings.Contains(string(config.Data), "foobar") {
		t.Fatal("data doesn't contain foobar step")
	}
	t.Log(config.Data)
}
