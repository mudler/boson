package main

import (
	. "github.com/mattn/go-getopt"
	_ "github.com/mudler/boson/processor/preprocessor"
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

const tmpdir = "/var/tmp/boson/"

func main() {

	var c int
	var configurationFile string
	//var logFile string
	OptErr = 0
	for {
		if c = Getopt("c:l:h"); c == EOF {
			break
		}
		switch c {
		case 'c':
			configurationFile = OptArg
			//	case 'l':
			//		logFile = OptArg
		case 'h':
			//		println("usage: " + os.Args[0] + " [-c my-boson-file.yaml|-l logfile|-h]")
			os.Exit(1)
		}
	}

	if configurationFile == "" {
		panic("I can't work without a configuration file")
	}

	//   backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	// Only errors and more severe messages should be sent to backend1
	//   backend1Leveled := logging.AddModuleLevel(backend1)
	//   backend1Leveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(backend2Formatter)

	log.Info("Loading config")

	config, err := utils.LoadConfig(configurationFile)

	if err != nil {
		panic(err)
	}

	// Bootstrapper for plugins
	for i, _ := range plugin_registry.Preprocessors {
		log.Info("Preprocessor found:" + i)
		plugin_registry.Preprocessors[i].OnStart()
	}

	// var build jdb.Build
	// build.Id = "2"
	// build.Passed = true
	// build.Commit = "test"
	// client.SaveBuild(build)
	// os.Exit(1)

	ticker := time.NewTicker(time.Second * time.Duration(config.PollTime))
	//go func() {
	os.MkdirAll(tmpdir, 666)
	workdir := tmpdir + config.RepositoryStripped
	client := jdb.NewDB("./" + configurationFile + ".db")
	for t := range ticker.C {
		log.Debug("Tick at", t)
		log.Debug("Cloning " + config.Repository + " to " + workdir)
		if ok, _ := utils.Exists(workdir); ok == true { //if already exists, using fetch && reset
			head := utils.GitHead(workdir)
			log.Info(utils.Git([]string{"fetch", "--all"}, workdir))
			log.Info(utils.Git([]string{"reset", "--hard", "origin/master"}, workdir))
			currentbuild, _ := client.GetBuild("LATEST_PASSED")
			log.Info("Head now is at " + head)
			if head == currentbuild.Commit {
				log.Info("nothing to do")
				continue
			}

			ContainerArgs, ContainerVolumes := plugin_registry.Preprocessors[config.PreProcessor].Process(workdir, &config, client)

			if ok, _ := utils.ContainerDeploy(&config, ContainerArgs, ContainerVolumes); ok == true {
				build := jdb.Build{Id: "LATEST_PASSED", Passed: true, Commit: head}
				client.SaveBuild(build)
			}
			//	deploy(&config, []string{"app-text/tree"})
		} else { //otherwise simply clone the repo22
			log.Info(utils.Git([]string{"clone", config.Repository, workdir}, tmpdir))
		}
	}
	//}()
	//time.Sleep(time.Millisecond * 1500)
	//ticker.Stop()
	//fmt.Println("Ticker stopped")

}
