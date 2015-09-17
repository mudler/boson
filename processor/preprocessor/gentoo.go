package gentoo

import (
	"github.com/mudler/boson/jdb"
	"github.com/mudler/boson/shared/registry"
	"github.com/mudler/boson/shared/utils"
	"github.com/op/go-logging"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var log = logging.MustGetLogger("boson")

type Gentoo struct{}

func (g *Gentoo) Process(workdir string, config *utils.Config, db *jdb.BuildClient) ([]string, []string) { //returns args and volumes to mount

	var ebuilds []string
	var volumes []string
	head := utils.GitHead(workdir)
	build, err := db.GetBuild("LATEST_PASSED")
	if err != nil {
		log.Debug("Database returned no result")
		log.Debug(err.Error())
		//return []string{}, map[string]string{}
	}
	log.Info("Commit: " + build.Commit)
	diffs, _ := utils.Git([]string{"diff", build.Commit, head, "--name-only"}, workdir)
	files := strings.Split(diffs, "\n")

	for _, v := range files {
		filename, _ := filepath.Abs(workdir + "/" + v)
		_, err := ioutil.ReadFile(filename)
		if err == nil {
			if strings.Contains(v, ".ebuild") { // We just need ebuilds
				eparts := strings.Split(strings.Replace(v, ".ebuild", "", -1), "/")
				ebuilds = append(ebuilds, "="+eparts[0]+"/"+eparts[2])
			}
		}
	}

	for _, e := range ebuilds {
		log.Debug(e)
	}

	volumes = append(volumes, workdir+":/usr/local/portage:ro") //my volume dir to mount
	if config.SeparateArtifacts == true {
		//in such case, we explicitally want to separate each artifacts directory (in gentoo case , our artifact is /usr/portage/packages)
		volumes = append(volumes, config.Artifacts+"/"+head+":/usr/portage/packages")
	} else {
		volumes = append(volumes, config.Artifacts+":/usr/portage/packages")
	}
	return ebuilds, volumes
}

func (g *Gentoo) OnStart() {
	log.Info("Gentoo Preprocessor")
}

func init() {
	plugin_registry.RegisterPreprocessor(&Gentoo{})
}
