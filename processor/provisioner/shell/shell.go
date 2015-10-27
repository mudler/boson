package shell

import (
	"github.com/mudler/boson/jdb"
	"github.com/mudler/boson/shared/registry"
	"github.com/mudler/boson/shared/utils"

	"github.com/op/go-logging"
)

// Shell construct the container arguments from the boson file
type Shell struct{}

var log = logging.MustGetLogger("boson")

// OnStart is the Shell entrypoint
func (s *Shell) OnStart() {
	log.Info("Shell Provisioner")
}

// Process builds a list of packages from the boson file
func (s *Shell) Process(workdir string, config *utils.Config, db *jdb.BuildClient) ([]string, []string) {

	return config.Provisioner["shell.Shell"], []string{}
}

func init() {
	pluginregistry.RegisterProvisioner(&Shell{})
}
