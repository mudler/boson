package main

import (
	. "github.com/mattn/go-getopt"

	"fmt"
	"os"
	"time"

	"github.com/mudler/boson/bosons"
	"github.com/mudler/boson/jdb"
	_ "github.com/mudler/boson/processor/preprocessor/gentoo"
	_ "github.com/mudler/boson/processor/provisioner/shell"
	"github.com/mudler/boson/shared/utils"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("boson")

var debugformat = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortpkg}.%{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
)
var normalformat = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
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

	if configurationFile == "1" {
		fmt.Println("I can't work without a configuration file")
		os.Exit(1)
	}

	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, debugformat)
	backenddebug2Formatter := logging.NewBackendFormatter(backend2, normalformat)

	if os.Getenv("DEBUG") == "1" {
		logging.SetBackend(backenddebug2Formatter)

	} else {
		logging.SetBackend(backend2Formatter)
	}

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
	builder := boson.NewBuilder(&config)
	if ok, _ := utils.Exists(config.WorkDir); ok == true { //if already exists, using fetch && reset
		utils.GitAlignToUpstream(config.WorkDir)
	} else { //otherwise simply clone the repo
		log.Info(utils.Git([]string{"clone", config.Repository, config.WorkDir}, config.TmpDir))
	}
	if _, ok := boson.Preprocessors[config.PreProcessor]; ok {

		if os.Getenv("BOSON_FROM") != "" && os.Getenv("BOSON_TO") != "" {
			builder.Run(builder.NewBuild(os.Getenv("BOSON_TO"), os.Getenv("BOSON_FROM")))
		} else if os.Getenv("BOSON_FROM") != "" {

			head := utils.GitHead(config.WorkDir)
			ok, err := builder.Run(builder.NewBuild(head, os.Getenv("BOSON_FROM")))
			if ok == true && err == nil {
				os.Exit(0)
			} else {
				os.Exit(42)
			}

		} else { //defaulting to ticker mode

			ticker := time.NewTicker(time.Second * time.Duration(config.PollTime))
			for _ = range ticker.C {
				log.Debug(" Cloning " + config.Repository + " to " + config.WorkDir)
				head := utils.GitHead(config.WorkDir)

				if ok, _ := utils.Exists(config.WorkDir); ok == true { //if already exists, using fetch && reset
					utils.GitAlignToUpstream(config.WorkDir)
				} else { //otherwise simply clone the repo
					log.Info(utils.Git([]string{"clone", config.Repository, config.WorkDir}, config.TmpDir))
				}

				lastbuild, _ := client.GetBuild("LATEST_PASSED")
				log.Info("Head now is at " + head)
				if head == lastbuild.Commit {
					log.Info("nothing to do")
					continue
				}
				build := builder.NewBuild(head, lastbuild.Commit)
				if ok, _ = builder.Run(build); ok == true { //Save the build status to id
					result := jdb.Build{Id: "LATEST_PASSED", Passed: true, Commit: build.Commit}
					client.SaveBuild(result)
					result = jdb.Build{Id: build.Commit, Passed: true, Commit: build.PrevCommit}
					client.SaveBuild(result)
				} else {
					result := jdb.Build{Id: "LATEST_PASSED", Passed: false, Commit: build.Commit}
					client.SaveBuild(result)
					result = jdb.Build{Id: build.Commit, Passed: false, Commit: build.PrevCommit}
					client.SaveBuild(result)
				}

			}
		}
	} else { //Provisioning
		if builder.Provision(builder.NewBuild("", "")) == true {
			os.Exit(0)
		} else {
			os.Exit(42)
		}
	}

}
