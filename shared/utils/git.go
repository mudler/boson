package utils

import (
	"os"
	"os/exec"
	"strings"
)

// Git executes a git command with the given args as a []string, outputs as a string
func Git(cmdArgs []string, dir string) (string, error) {
	var (
		cmdOut []byte
		err    error
	)
	cwd, _ := os.Getwd()
	os.Chdir(dir)
	cmdName := "git"
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		log.Error("There was an error running git command: ", err)
		log.Debug(strings.Join(cmdArgs, " "))
		log.Error(string(cmdOut))
	}
	os.Chdir(cwd)
	result := string(cmdOut)
	return strings.TrimSpace(result), err
}

// GitAlignToUpstream executes a fetch --all and reset --hard to origin/master on the given git repository
func GitAlignToUpstream(workdir string) {
	log.Info(Git([]string{"fetch", "--all"}, workdir))
	log.Info(Git([]string{"reset", "--hard", "origin/master"}, workdir))
}

func GitPrevCommit(workdir string) (string, error) {
	result, err := Git([]string{"log", "-2", `--pretty=format:"%h"`}, workdir)
	temp := strings.Split(result, "\n")
	return temp[1], err
}

// GitHead returns the Head of the given repository
func GitHead(workdir string) string {
	head, _ := Git([]string{"rev-parse", "HEAD"}, workdir)
	return head
}
