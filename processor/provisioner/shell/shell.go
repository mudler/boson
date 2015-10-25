package shell

import (
	"github.com/mudler/boson/jdb"
	"github.com/mudler/boson/shared/registry"
	"github.com/mudler/boson/shared/utils"

	"github.com/op/go-logging"
)

type Shell struct{}

var log = logging.MustGetLogger("boson")

func (s *Shell) OnStart() {
	log.Info("Shell Provisioner")
}
func (s *Shell) Process(workdir string, config *utils.Config, db *jdb.BuildClient) ([]string, []string) {

	return config.Provisioner["shell.Shell"], []string{}
}

func init() {
	plugin_registry.RegisterProvisioner(&Shell{})
}
