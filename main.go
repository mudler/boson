package main

import (
	. "github.com/mattn/go-getopt"

	_ "github.com/mudler/boson/processor/preprocessor/gentoo"
	_ "github.com/mudler/boson/processor/provisioner/shell"

	"github.com/mudler/boson/shared/registry"

	"github.com/mudler/boson/jdb"
	"os"
	"time"

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
		panic("I can't work without a configuration file")
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

	for i, _ := range plugin_registry.Preprocessors {
		log.Info("\t *" + i)
		plugin_registry.Preprocessors[i].OnStart()
	}

	os.MkdirAll(config.TmpDir, 666)
	workdir := config.TmpDir + config.RepositoryStripped
	client := jdb.NewDB("./" + configurationFile + ".db")

	if _, ok := plugin_registry.Preprocessors[config.PreProcessor]; ok {
		ticker := time.NewTicker(time.Second * time.Duration(config.PollTime))
		for range ticker.C {
			log.Debug("Cloning " + config.Repository + " to " + workdir)
			if ok, _ := utils.Exists(workdir); ok == true { //if already exists, using fetch && reset
				utils.GitAlignToUpstream(workdir)
				currentbuild, _ := client.GetBuild("LATEST_PASSED")
				head := utils.GitHead(workdir)
				log.Info("Head now is at " + head)
				if head == currentbuild.Commit {
					log.Info("nothing to do")
					continue
				}

				ContainerArgs, ContainerVolumes := plugin_registry.Preprocessors[config.PreProcessor].Process(workdir, &config, client)

				if ok, _ := utils.ContainerDeploy(&config, ContainerArgs, ContainerVolumes, head); ok == true {
					build := jdb.Build{Id: "LATEST_PASSED", Passed: true, Commit: head}
					client.SaveBuild(build)
					build = jdb.Build{Id: head, Passed: true, Commit: currentbuild.Commit}
					client.SaveBuild(build)
				} else {
					build := jdb.Build{Id: "LATEST_PASSED", Passed: false, Commit: head}
					client.SaveBuild(build)
					build = jdb.Build{Id: head, Passed: false, Commit: currentbuild.Commit}
					client.SaveBuild(build)
				}

				for i, _ := range plugin_registry.Postprocessors {
					log.Info("Postprocessor found:" + i)
					plugin_registry.Postprocessors[i].Process(workdir, &config, client)
				}
			} else { //otherwise simply clone the repo
				log.Info(utils.Git([]string{"clone", config.Repository, workdir}, config.TmpDir))
			}
		}
	} else { //Provisioning
		for i, _ := range config.Provisioner {
			log.Info("\t - " + i)
			plugin_registry.Provisioners[i].OnStart()
			ContainerArgs, ContainerVolumes := plugin_registry.Provisioners[i].Process(workdir, &config, client)

			if ok, _ := utils.ContainerDeploy(&config, ContainerArgs, ContainerVolumes, "LATEST-PROVISIONED"); ok == true {
				log.Info("All done")
				os.Exit(0)
			} else {
				log.Error("Build failed")

				os.Exit(100)
			}

		}
	}

}
