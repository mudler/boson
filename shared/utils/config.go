package utils

import (
	"github.com/op/go-logging"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v2"
)

var log = logging.MustGetLogger("boson")

type Config struct {
	// Firewall_network_rules map[string]Options `yaml:"repository"`
	Repository         string `yaml:"repository"`
	RepositoryStripped string
	DockerImage        string `yaml:"docker_image"`
	PreProcessor       string `yaml:"preprocessor"`
	PostProcessor      string `yaml:"postprocessor"`
	PollTime           int    `yaml:"polltime"`
	Artifacts          string `yaml:"artifacts_dir"`
	SeparateArtifacts  bool   `yaml:"separate_artifacts"`
	LogDir             string `yaml:"log_dir"`
	LogPerm            int    `yaml:"logfile_perm"`
}

//type Options struct {
//   Src string
//   Dst string
//}

func LoadConfig(f string) (Config, error) {

	filename, _ := filepath.Abs(f)
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var config Config
	config.SeparateArtifacts = false
	config.PollTime = 5
	config.LogPerm = int(777)
	err = yaml.Unmarshal(yamlFile, &config)

	r, _ := regexp.Compile(`^.*?\/\/`)
	config.RepositoryStripped = r.ReplaceAllString(config.Repository, "")

	if config.Artifacts == "" {
		log.Fatal("You need to specify 'artifacts_dir'")
	}
	if config.PreProcessor == "" {
		log.Fatal("You need to specify a preprocessor 'preprocessors'")
	}
	if config.Repository == "" {
		log.Fatal("You need to specify a repository 'repository'")
	}
	if config.DockerImage == "" {
		log.Fatal("You need to specify a Docker image 'docker_image'")
	}
	if config.LogDir == "" {
		log.Fatal("You need to specify a Log directory 'log_dir'")
	}
	log.Info("GIT Repository: %#v\n", config.Repository)
	log.Info("Docker Image: %#v\n", config.DockerImage)
	log.Info("Artifacts directory: %#v\n", config.Artifacts)
	log.Info("Separate Artifacts by commit: %#t\n", config.SeparateArtifacts)

	log.Info("PreProcessor: %#v\n", config.PreProcessor)
	log.Info("Log Directory: %#v\n", config.LogDir)
	log.Info("Log Permissions: %#d\n", config.LogPerm)
	log.Info("Poll Time: %#v\n", config.PollTime)

	return config, err
}
