package gentoo

import (
	"github.com/mudler/boson/jdb"
	"github.com/mudler/boson/shared/registry"
	"github.com/mudler/boson/shared/utils"
	"github.com/op/go-logging"
	"strings"
)

var log = logging.MustGetLogger("boson")

type Gentoo struct{}

func (g *Gentoo) Process(workdir string, config *utils.Config, db *jdb.BuildClient) ([]string, []string) { //returns args and volumes to mount

	var ebuilds []string
	var volumes []string
	build, err := db.GetBuild("LATEST_PASSED")
	if err != nil {
		log.Debug("Database returned no result")
		log.Debug(err.Error())
		//return []string{}, map[string]string{}
	}
	log.Info("Commit: " + build.Commit)
	diffs, _ := utils.Git([]string{"diff", build.Commit, utils.GitHead(workdir), "--name-only"}, workdir)
	files := strings.Split(diffs, "\n")

	for _, v := range files {
		if strings.Contains(v, ".ebuild") { // We just need ebuilds
			eparts := strings.Split(strings.Replace(v, ".ebuild", "", -1), "/")
			ebuilds = append(ebuilds, "="+eparts[0]+"/"+eparts[2])
		}
	}

	for _, e := range ebuilds {
		log.Debug(e)
	}

	volumes = append(volumes, workdir+":/usr/local/portage:ro") //my volume dir to mount
	volumes = append(volumes, config.Artifacts+":/usr/portage/packages")
	return ebuilds, volumes
}

func (g *Gentoo) OnStart() {
	log.Info("Gentoo Preprocessor")
}

func init() {
	plugin_registry.RegisterPreprocessor(&Gentoo{})
}
