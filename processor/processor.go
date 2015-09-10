package processor

import (
	"github.com/mudler/boson/shared/utils"
)

type Processor interface {
	Process(string, *utils.Config) []string // processor gets the workdir and the config file
	OnStart()
}
