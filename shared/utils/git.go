package utils

import (
	"os"
	"os/exec"
	"strings"
)

func Git(cmdArgs []string, dir string) (string, error) {
	var (
		cmdOut []byte
		err    error
	)
	cwd, _ := os.Getwd()
	os.Chdir(dir)
	cmdName := "git"
	//cmdArgs := []string{"rev-parse", "--verify", "HEAD"}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		log.Error("There was an error running git command: ", err)
		log.Debug(strings.Join(cmdArgs, " "))
		log.Error(string(cmdOut))
	}
	os.Chdir(cwd)
	result := string(cmdOut)
	return strings.TrimSpace(result), err
}

func GitAlignToUpstream(workdir string) {
	log.Info(Git([]string{"fetch", "--all"}, workdir))
	log.Info(Git([]string{"reset", "--hard", "origin/master"}, workdir))
}

func GitHead(workdir string) string {
	head, _ := Git([]string{"rev-parse", "HEAD"}, workdir)
	return head
}
