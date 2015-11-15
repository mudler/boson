package gentoo

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mudler/boson/bosons"
	"github.com/mudler/boson/shared/utils"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("boson")

const artifactsdir = "/usr/portage/packages"

// Gentoo is the Preprocessor that detects Gentoo ebuilds between commits
type Gentoo struct{}

// Process builds a list of packages to emerge between commits
func (g *Gentoo) Process(build *boson.Build) ([]string, []string) { //returns args and volumes to mount

	workdir := build.Config.WorkDir
	config := build.Config
	var ebuilds []string
	var volumes []string
	log.Info("Commit: " + build.Commit)
	diffs, _ := utils.Git([]string{"diff", build.PrevCommit, build.Commit, "--name-only"}, workdir)
	files := strings.Split(diffs, "\n")

	for _, v := range files {
		filename, _ := filepath.Abs(workdir + "/" + v)
		_, err := ioutil.ReadFile(filename)
		if err == nil {
			if strings.Contains(v, ".ebuild") { // We just need ebuilds
				eparts := strings.Split(strings.Replace(v, ".ebuild", "", -1), "/")
				ebuilds = append(ebuilds, "="+eparts[0]+"/"+eparts[2])
				build.Extras = append(build.Extras, "="+eparts[0]+"/"+eparts[2])
			}
		}
	}

	for _, e := range ebuilds {
		log.Debug(e)
	}

	volumes = append(volumes, workdir+":/usr/local/portage:ro") //my volume dir to mount
	if config.SeparateArtifacts == true {
		//in such case, we explicitally want to separate each artifacts directory (in gentoo , our artifact is /usr/portage/packages)
		volumes = append(volumes, config.Artifacts+"/"+build.Commit+":"+artifactsdir)
	} else {
		volumes = append(volumes, config.Artifacts+":"+artifactsdir)
	}
	return ebuilds, volumes
}

// OnStart on Gentoo Preprocessor
func (g *Gentoo) OnStart() {
	log.Info("Gentoo Preprocessor")
}

func init() {
	boson.RegisterPreprocessor(&Gentoo{})
}
