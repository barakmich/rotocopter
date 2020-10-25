package main

import (
	"net/http"

	"github.com/barakmich/rotocopter/roto"
	"github.com/drone/drone-go/plugin/config"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// spec provides the plugin settings.
type spec struct {
	Bind   string `envconfig:"DRONE_BIND"`
	Debug  bool   `envconfig:"DRONE_DEBUG"`
	Secret string `envconfig:"DRONE_SECRET"`

	GitYamlRepo     string `envconfig:"DRONE_YAML_GIT_REPO"`
	GitHTTPUser     string `envconfig:"DRONE_YAML_GIT_USER"`
	GitHTTPPassword string `envconfig:"DRONE_YAML_GIT_PASSWORD"`
}

func main() {
	spec := new(spec)
	err := envconfig.Process("", spec)
	if err != nil {
		logrus.Fatal(err)
	}

	if spec.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if spec.Secret == "" {
		logrus.Fatalln("missing secret key")
	}
	if spec.Bind == "" {
		spec.Bind = ":3000"
	}
	if spec.GitYamlRepo == "" {
		logrus.Fatalln("need to reference a git repo")
	}

	conf := roto.Config{
		GitYamlRepo:     spec.GitYamlRepo,
		GitHTTPUser:     spec.GitHTTPUser,
		GitHTTPPassword: spec.GitHTTPPassword,
	}

	handler := config.Handler(
		roto.New(conf),
		spec.Secret,
		logrus.StandardLogger(),
	)

	logrus.Infof("server listening on address %s", spec.Bind)

	http.Handle("/", handler)
	logrus.Fatal(http.ListenAndServe(spec.Bind, nil))
}
