package jdb

import (
	"io/ioutil"
	"log"
)

func NewDB(db string) *BuildClient {
	if _, err := ioutil.ReadFile(db); err != nil {
		str := "{}"
		if err = ioutil.WriteFile(db, []byte(str), 0644); err != nil {
			log.Fatal(err)
		}
	}

	// create channel to communicate over
	jobs := make(chan Job)

	// start watching jobs channel for work
	go ProcessJobs(jobs, db)

	// create client for submitting jobs / providing interface to db
	// create dependencies
	client := &BuildClient{Jobs: jobs}

	return client

}
