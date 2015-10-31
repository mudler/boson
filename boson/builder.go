package boson

import (
	"github.com/mudler/boson/jdb"
	"github.com/mudler/boson/shared/utils"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("boson")

// Builder represent the engine that runs the builds. It has a *utils.Config and a jdb.BuildClient
type Builder struct {
	Config *utils.Config
	db     *jdb.BuildClient
}

// NewBuilder returns a new *boson.Builder
func NewBuilder(config *utils.Config, db *jdb.BuildClient) *Builder {
	return &Builder{Config: config, db: db}
}

// NewBuild returns a new *boson.Build
func (b *Builder) NewBuild(Commit, PrevCommit string) *Build {
	return &Build{Config: b.Config, Commit: Commit, PrevCommit: PrevCommit}
}

// Run runs the given *boson.Build
func (b *Builder) Run(build *Build) {
	ContainerArgs, ContainerVolumes := Preprocessors[build.Config.PreProcessor].Process(build, b.db)
	log.Info(">> Running build on " + b.Config.DockerImage)
	//Save state
	if ok, _ := utils.ContainerDeploy(build.Config, ContainerArgs, ContainerVolumes, build.Commit); ok == true {
		result := jdb.Build{Id: "LATEST_PASSED", Passed: true, Commit: build.Commit}
		b.db.SaveBuild(result)
		result = jdb.Build{Id: build.Commit, Passed: true, Commit: build.PrevCommit}
		b.db.SaveBuild(result)
	} else {
		result := jdb.Build{Id: "LATEST_PASSED", Passed: false, Commit: build.Commit}
		b.db.SaveBuild(result)
		result = jdb.Build{Id: build.Commit, Passed: false, Commit: build.PrevCommit}
		b.db.SaveBuild(result)
	}

	for i := range Postprocessors {
		Postprocessors[i].Process(build, b.db)
	}
}

// Provision provisions the given *boson.Build
func (b *Builder) Provision(build *Build) bool {
	ret := true
	for i := range build.Config.Provisioner {
		Provisioners[i].OnStart()
		ContainerArgs, ContainerVolumes := Provisioners[i].Process(build, b.db)

		if ok, _ := utils.ContainerDeploy(build.Config, ContainerArgs, ContainerVolumes, "LATEST-PROVISIONED"); ok == false {
			ret = false
		}

	}
	return ret
}
