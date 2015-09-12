package processor

import (
	"github.com/mudler/boson/jdb"
	"github.com/mudler/boson/shared/utils"
)

type Processor interface {
	Process(string, *utils.Config, *jdb.BuildClient) ([]string, map[string]string) // processor gets the workdir and the config file
	OnStart()
}
