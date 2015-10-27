package processor

import (
	"github.com/mudler/boson/jdb"
	"github.com/mudler/boson/shared/utils"
)

// Processor is the interface for preprocessors, postprocessors and provisioners
type Processor interface {
	Process(string, *utils.Config, *jdb.BuildClient) ([]string, []string) // processor gets the workdir and the config file
	OnStart()
}
