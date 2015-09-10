package jdb

type BuildClient struct {
	Jobs chan Job
}

func (c *BuildClient) SaveBuild(build Build) (Build, error) {
	job := NewSaveBuildJob(build)
	c.Jobs <- job

	if err := <-job.ExitChan(); err != nil {
		return Build{}, err
	}
	return <-job.saved, nil
}

func (c *BuildClient) GetBuilds() ([]Build, error) {
	arr := make([]Build, 0)

	builds, err := c.getBuildHash()
	if err != nil {
		return arr, err
	}

	for _, value := range builds {
		arr = append(arr, value)
	}
	return arr, nil
}

func (c *BuildClient) GetBuild(id string) (Build, error) {
	builds, err := c.getBuildHash()
	if err != nil {
		return Build{}, err
	}
	return builds[id], nil
}

func (c *BuildClient) DeleteBuild(id string) error {
	job := NewDeleteBuildJob(id)
	c.Jobs <- job

	if err := <-job.ExitChan(); err != nil {
		return err
	}
	return nil
}

func (c *BuildClient) getBuildHash() (map[string]Build, error) {
	job := NewReadBuildsJob()
	c.Jobs <- job

	if err := <-job.ExitChan(); err != nil {
		return make(map[string]Build, 0), err
	}
	return <-job.builds, nil
}
