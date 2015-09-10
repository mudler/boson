package jdb

// Job to read all builds from the database
type ReadBuildsJob struct {
	builds   chan map[string]Build
	exitChan chan error
}

func NewReadBuildsJob() *ReadBuildsJob {
	return &ReadBuildsJob{
		builds:   make(chan map[string]Build, 1),
		exitChan: make(chan error, 1),
	}
}
func (j ReadBuildsJob) ExitChan() chan error {
	return j.exitChan
}
func (j ReadBuildsJob) Run(builds map[string]Build) (map[string]Build, error) {
	j.builds <- builds

	return nil, nil
}
