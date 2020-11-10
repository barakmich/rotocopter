package roto

import (
	"io/ioutil"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/config"
	"github.com/go-git/go-git/v5"
	"github.com/sirupsen/logrus"
)

func getDroneYaml(val string, wt *git.Worktree, req *config.Request) (*drone.Config, error) {
	if strings.HasSuffix(val, ".yaml") || strings.HasSuffix(val, ".yml") {
		return getDroneYamlFromFile(val, wt)
	}
	if i := strings.Index(val, ":"); i >= 0 {
		file := val[:i]
		switch {
		case strings.HasSuffix(file, ".py"):
			fallthrough
		case strings.HasSuffix(file, ".bzl"):
			fallthrough
		case strings.HasSuffix(file, ".star"):
			return getDroneYamlFromStarlark(val, wt, req)
		}
	}
	return getDroneYamlFromFile(val, wt)
}

func getDroneYamlFromFile(v string, wt *git.Worktree) (*drone.Config, error) {
	yamlFile, err := wt.Filesystem.Open(v)
	if err != nil {
		logrus.Error("can't open filename", v)
		return nil, err
	}
	data, err := ioutil.ReadAll(yamlFile)
	if err != nil {
		return nil, err
	}

	conf := &drone.Config{
		Data: string(data),
	}
	// return nil and Drone will fallback to
	// the standard behavior for getting the
	// configuration file.
	return conf, nil
}

func getDroneYamlFromStarlark(v string, wt *git.Worktree, req *config.Request) (*drone.Config, error) {

}
