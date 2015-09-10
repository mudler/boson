package utils

import (
	"os"
	"os/exec"
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
		log.Error(string(cmdOut))
		return "", err
	}
	os.Chdir(cwd)
	result := string(cmdOut)
	return result, err
}

func GitHead(workdir string) string {
	head, _ := Git([]string{"rev-parse", "HEAD"}, workdir)
	return head
}
