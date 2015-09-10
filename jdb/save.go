package jdb

import (
	"crypto/rand"
	"fmt"
	"io"
)

// Job to add a Build to the database
type SaveBuildJob struct {
	toSave   Build
	saved    chan Build
	exitChan chan error
}

func NewSaveBuildJob(build Build) *SaveBuildJob {
	return &SaveBuildJob{
		toSave:   build,
		saved:    make(chan Build, 1),
		exitChan: make(chan error, 1),
	}
}
func (j SaveBuildJob) ExitChan() chan error {
	return j.exitChan
}
func (j SaveBuildJob) Run(builds map[string]Build) (map[string]Build, error) {
	var build Build
	if j.toSave.Id == "" {
		id, err := newUUID()
		if err != nil {
			return nil, err
		}
		build = Build{Id: id, Passed: j.toSave.Passed, Commit: j.toSave.Commit}
	} else {
		build = j.toSave
	}
	builds[build.Id] = build

	j.saved <- build
	return builds, nil
}

// Generate a uuid to use as a unique identifier for each Build
// http://play.golang.org/p/4FkNSiUDMg
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
