package jdb

// Job to delete a Build from the database
type DeleteBuildJob struct {
	toDelete string
	exitChan chan error
}

func NewDeleteBuildJob(id string) *DeleteBuildJob {
	return &DeleteBuildJob{
		toDelete: id,
		exitChan: make(chan error, 1),
	}
}
func (j DeleteBuildJob) ExitChan() chan error {
	return j.exitChan
}
func (j DeleteBuildJob) Run(builds map[string]Build) (map[string]Build, error) {
	delete(builds, j.toDelete)
	return builds, nil
}
