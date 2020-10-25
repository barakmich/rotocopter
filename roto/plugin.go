package roto

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/config"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/plumbing/transport"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/sirupsen/logrus"
)

// New returns a new config plugin.
func New(config Config) config.Plugin {
	return newRotoPlugin(config)
}

func newRotoPlugin(config Config) *rotoPlugin {
	plugin := &rotoPlugin{}

	logrus.Infof("cloning %s", config.GitYamlRepo)

	if config.GitHTTPUser != "" {
		plugin.auth = &http.BasicAuth{
			Username: config.GitHTTPUser,
			Password: config.GitHTTPPassword,
		}
	}

	fs := memfs.New()
	storer := memory.NewStorage()

	repo, err := git.Clone(storer, fs, &git.CloneOptions{
		Auth: plugin.auth,
		URL:  config.GitYamlRepo,
	})
	if err != nil {
		logrus.Fatal("couldn't clone repo:", err)
	}
	plugin.repo = repo
	err = plugin.updateIndex()
	if err != nil {
		logrus.Fatal("couldn't update initial index", err)
	}
	return plugin
}

type rotoPlugin struct {
	auth  transport.AuthMethod
	repo  *git.Repository
	index map[string]string
}

func (r *rotoPlugin) updateWorktree() error {
	wt, err := r.repo.Worktree()
	if err != nil {
		return err
	}

	logrus.Info("pulling tree")
	err = wt.Pull(&git.PullOptions{
		Auth: r.auth,
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			logrus.Info("up to date")
			return nil
		} else {
			logrus.Errorln("couldn't pull:", err)
			return err
		}
	}
	return r.updateIndex()
}

func (r *rotoPlugin) updateIndex() error {
	wt, err := r.repo.Worktree()
	if err != nil {
		return err
	}
	f, err := wt.Filesystem.Open("index.json")
	if err != nil {
		return err
	}
	index := make(map[string]string)
	err = json.NewDecoder(f).Decode(&index)
	if err != nil {
		return err
	}
	r.index = index
	for reposlug := range r.index {
		logrus.Infof("Repo %s", reposlug)
	}
	return nil
}

func (r *rotoPlugin) Find(ctx context.Context, req *config.Request) (*drone.Config, error) {
	err := r.updateWorktree()
	if err != nil {
		return nil, err
	}
	wt, err := r.repo.Worktree()
	if err != nil {
		logrus.Error("can't get worktree:", err)
		return nil, err
	}
	v, ok := r.index[req.Repo.Slug]
	if !ok {
		logrus.Info(req.Repo.Slug, " doesn't exist in index, skipping")
		return nil, nil
	}
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
