package boson

import (
	"github.com/mudler/boson/jdb"
)

// Processor is the interface for preprocessors, postprocessors and provisioners
type Processor interface {
	Process(*Build, *jdb.BuildClient) ([]string, []string) // processor gets the workdir and the config file
	OnStart()
}
