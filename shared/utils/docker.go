package utils

import (
	"github.com/samalba/dockerclient"
)

// Callback used to listen to Docker's events
func eventCallback(event *dockerclient.Event, ec chan error, args ...interface{}) {
	log.Info("Docker instance: %#v\n", *event)
}

func ContainerDeploy(config *Config, args []string, f func(string)) (string, error) {

	docker, _ := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)

	DockerImage := config.DockerImage
	builderConfig := &dockerclient.ContainerConfig{
		Image: DockerImage,
		//Entrypoint: []string{"/bin/bash"},
		Cmd: args,
		//  Tty:        true,
		//OpenStdin:  true,
		//  HostConfig: dockerclient.HostConfig{
		//      RestartPolicy: dockerclient.RestartPolicy{
		//          Name:              "always",
		//          MaximumRetryCount: 0,
		//      },
		//  },
	}

	// pull images
	log.Info("Pulling image: %s\n", DockerImage)
	if err := docker.PullImage(DockerImage, nil); err != nil {
		log.Error("error pulling %s image: %s\n", DockerImage, err)
		return "", err
	}

	// start the show
	builderId, err := docker.CreateContainer(builderConfig, "boson-builder")
	if err != nil {
		log.Error("error creating %s: %s\n", DockerImage, err)
		return "", err
	}
	if err := docker.StartContainer(builderId, &builderConfig.HostConfig); err != nil {
		log.Error("error starting %s: %s\n", DockerImage, err)
		return "", err
	}
	docker.StartMonitorEvents(eventCallback, nil)

	// Remove handling -> gofunc to wait && delete
	// Valid case.
	//go func() {

	/// MEGLIO SENZA GOFUNC? in teoria dovrei attendere tra un job e un altro
	select {
	case wr := <-docker.Wait(builderId):
		if wr.ExitCode == int(0) {
			// success, call the callback to save to db last valid commit
			docker.RemoveContainer(builderId, true, false)
		}
	}
	//}()
	log.Info("boson Stack started successfully")
	return builderId, err
}
