package gentoo

import (
	"github.com/mudler/boson/shared/registry"
	"github.com/mudler/boson/shared/utils"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("boson")

type Gentoo struct{}

func (g *Gentoo) Process(workdir string, config *utils.Config) []string {
	return []string{"albero"}
}

func (g *Gentoo) OnStart() {
	log.Info("Hi from Gentoo preprocessor")
}

func init() {
	plugin_registry.RegisterPreprocessor(&Gentoo{})
}
