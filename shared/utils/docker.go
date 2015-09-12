package utils

import (
	"github.com/samalba/dockerclient"
	"os"
)

// Callback used to listen to Docker's events
func eventCallback(event *dockerclient.Event, ec chan error, args ...interface{}) {
	log.Info("Docker instance: %#v\n", *event)
}

func ContainerDeploy(config *Config, args []string, volumes map[string]string) (bool, error) {

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
		return false, err
	}

	// start the show
	builderId, err := docker.CreateContainer(builderConfig, config.Repository)
	if err != nil {
		log.Error("error creating %s: %s\n", DockerImage, err)
		return false, err
	}
	var devices []dockerclient.DeviceMapping

	for i, c := range volumes {
		log.Info("Host path: " + c)
		log.Info("Container path:" + volumes[i])
		devices = append(devices, dockerclient.DeviceMapping{
			PathOnHost:      c,
			PathInContainer: volumes[i],
		})
	}
	os.Exit(1)

	hostConfig := &dockerclient.HostConfig{Devices: devices}
	if err := docker.StartContainer(builderId, hostConfig); err != nil {
		log.Error("error starting %s: %s\n", DockerImage, err)
		return false, err
	}
	docker.StartMonitorEvents(eventCallback, nil)

	// Remove handling -> gofunc to wait && delete
	// Valid case.
	//go func() {

	/// MEGLIO SENZA GOFUNC? in teoria dovrei attendere tra un job e un altro
	select {
	case wr := <-docker.Wait(builderId):
		docker.RemoveContainer(builderId, true, false)
		if wr.ExitCode == int(0) {
			// success, call the callback to save to db last valid commit
			return true, err
		} else {
			return false, err
		}
	}
}
