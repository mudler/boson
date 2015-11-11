package boson

import (
	"github.com/mudler/boson/shared/utils"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("boson")

// Builder represent the engine that runs the builds. It has a *utils.Config and a jdb.BuildClient
type Builder struct {
	Config *utils.Config
}

// NewBuilder returns a new *boson.Builder
func NewBuilder(config *utils.Config) *Builder {
	return &Builder{Config: config}
}

// NewBuild returns a new *boson.Build
func (b *Builder) NewBuild(Commit, PrevCommit string) *Build {
	return &Build{Config: b.Config, Commit: Commit, PrevCommit: PrevCommit}
}

// Run runs the given *boson.Build
func (b *Builder) Run(build *Build) (bool, error) {

	log.Info(">> Build " + build.PrevCommit + " --> " + build.Commit)

	ContainerArgs, ContainerVolumes := Preprocessors[build.Config.PreProcessor].Process(build)

	log.Info(">> Running build on " + b.Config.DockerImage)
	//Save state
	ok := false
	var err error
	ok, err = utils.ContainerDeploy(build.Config, ContainerArgs, ContainerVolumes, build.Commit)

	for i := range Postprocessors {
		Postprocessors[i].Process(build)
	}
	return ok, err
}

// Provision provisions the given *boson.Build
func (b *Builder) Provision(build *Build) bool {
	ret := true
	for i := range build.Config.Provisioner {
		Provisioners[i].OnStart()
		ContainerArgs, ContainerVolumes := Provisioners[i].Process(build)

		if ok, _ := utils.ContainerDeploy(build.Config, ContainerArgs, ContainerVolumes, "LATEST-PROVISIONED"); ok == false {
			ret = false
		}

	}
	return ret
}
