package main

import (
	. "github.com/mattn/go-getopt"

	"fmt"
	"os"
	"time"

	"github.com/mudler/boson/boson"
	"github.com/mudler/boson/jdb"
	_ "github.com/mudler/boson/processor/preprocessor/gentoo"
	_ "github.com/mudler/boson/processor/provisioner/shell"
	"github.com/mudler/boson/shared/utils"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("boson")
var format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortpkg}.%{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
)

func main() {

	var c int
	var configurationFile string
	OptErr = 0
	for {
		if c = Getopt("c:h"); c == EOF {
			break
		}
		switch c {
		case 'c':
			configurationFile = OptArg
		case 'h':
			println("usage: " + os.Args[0] + " [-c my-boson-file.yaml -h]")
			os.Exit(1)
		}
	}

	if configurationFile == "" {
		fmt.Println("I can't work without a configuration file")
		os.Exit(1)
	}

	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	logging.SetBackend(backend2Formatter)

	log.Info("Loading config")

	config, err := utils.LoadConfig(configurationFile)

	if err != nil {
		panic(err)
	}

	// Bootstrapper for plugins
	log.Info("Available preprocessors:")

	for i := range boson.Preprocessors {
		log.Info("\t *" + i)
		boson.Preprocessors[i].OnStart()
	}

	os.MkdirAll(config.TmpDir, 666)
	client := jdb.NewDB("./" + configurationFile + ".db")
	builder := boson.NewBuilder(&config, client)

	if _, ok := boson.Preprocessors[config.PreProcessor]; ok {
		ticker := time.NewTicker(time.Second * time.Duration(config.PollTime))
		for _ = range ticker.C {
			log.Debug(" Cloning " + config.Repository + " to " + config.WorkDir)
			if ok, _ := utils.Exists(config.WorkDir); ok == true { //if already exists, using fetch && reset
				utils.GitAlignToUpstream(config.WorkDir)
			} else { //otherwise simply clone the repo
				log.Info(utils.Git([]string{"clone", config.Repository, config.WorkDir}, config.TmpDir))
			}

			lastbuild, _ := client.GetBuild("LATEST_PASSED")
			head := utils.GitHead(config.WorkDir)
			log.Info("Head now is at " + head)
			if head == lastbuild.Commit {
				log.Info("nothing to do")
				continue
			}

			builder.Run(builder.NewBuild(head, lastbuild.Commit))

		}
	} else { //Provisioning
		if builder.Provision(builder.NewBuild("", "")) == true {
			os.Exit(0)
		} else {
			os.Exit(42)
		}
	}

}
