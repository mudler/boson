package shell

import (
	"github.com/mudler/boson/bosons"

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
func (s *Shell) Process(build *boson.Build) ([]string, []string) { //returns args and volumes to mount
	return build.Config.Provisioner["shell.Shell"], []string{}
}

func init() {
	boson.RegisterProvisioner(&Shell{})
}
