package jdb

import (
	"encoding/json"
	"io/ioutil"
)

type Job interface {
	ExitChan() chan error
	Run(builds map[string]Build) (map[string]Build, error)
}

func ProcessJobs(jobs chan Job, db string) {
	for {
		j := <-jobs

		builds := make(map[string]Build, 0)
		content, err := ioutil.ReadFile(db)
		if err == nil {
			if err = json.Unmarshal(content, &builds); err == nil {
				buildsMod, err := j.Run(builds)

				if err == nil && buildsMod != nil {
					b, err := json.Marshal(buildsMod)
					if err == nil {
						err = ioutil.WriteFile(db, b, 0644)
					}
				}
			}
		}

		j.ExitChan() <- err
	}
}
