package utils

import (
	//"github.com/samalba/dockerclient"
	"bufio"
	"github.com/fsouza/go-dockerclient"
	"io"
	"os"
)

func Attach(client *docker.Client, container *docker.Container) error {
	r, w := io.Pipe()
	options := docker.AttachToContainerOptions{
		Container:    container.ID,
		OutputStream: w,
		ErrorStream:  w,
		Stream:       true,
		Stdout:       true,
		Stderr:       true,
		Logs:         true,
	}
	log.Info("Attaching to container " + container.ID)
	go client.AttachToContainer(options)
	go func(reader io.Reader, container *docker.Container) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			log.Info("%s%s \n", container.ID, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Error("There was an error with the scanner in attached container", err)
		}
	}(r, container)
	return nil
}
func ContainerDeploy(config *Config, args []string, volumes []string, head string) (bool, error) {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	DockerImage := config.DockerImage
	log.Info("Pulling image: %s\n", DockerImage)

	if err := client.PullImage(docker.PullImageOptions{Repository: DockerImage}, docker.AuthConfiguration{}); err != nil {
		log.Error("error pulling %s image: %s\n", DockerImage, err)
		return false, err
	}

	container, err := client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: DockerImage,
			Cmd:   args,
		},
	})
	Attach(client, container)
	// Cleanup when done
	defer func() {
		client.RemoveContainer(docker.RemoveContainerOptions{
			ID:    container.ID,
			Force: true,
		})
	}()
	if err != nil {
		log.Error(err.Error())
	}
	log.Info("Starting container: " + container.ID)
	err = client.StartContainer(container.ID, &docker.HostConfig{Binds: volumes, LogConfig: docker.LogConfig{Type: "json-file"}})
	if err != nil {
		log.Error(err.Error())
		return false, err
	}
	status, err := client.WaitContainer(container.ID)
	container, err = client.InspectContainer(container.ID)
	if err != nil {
		log.Error(err.Error())
	}
	// update our container information
	logfile := config.LogDir + "/" + head + ".json"
	log.Debug("Copying " + container.LogPath + " to " + logfile)
	err = CopyFile(container.LogPath, logfile)
	os.Chmod(logfile, os.FileMode(config.LogPerm))
	if err != nil {
		log.Error(err.Error())
	}

	if status == int(0) {
		return true, err
	} else {
		return false, err
	}
}

// Callback used to listen to Docker's events
// func eventCallback(event *dockerclient.Event, ec chan error, args ...interface{}) {
// 	log.Info("Docker instance: %#v\n", *event)
// }

// func ContainerDeploy(config *Config, args []string, volumes map[string]string) (bool, error) {

// 	docker, _ := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)

// 	DockerImage := config.DockerImage
// 	builderConfig := &dockerclient.ContainerConfig{
// 		Image: DockerImage,
// 		//Entrypoint: []string{"/bin/bash"},
// 		Cmd: args,
// 		//  Tty:        true,
// 		//OpenStdin:  true,
// 		//  HostConfig: dockerclient.HostConfig{
// 		//      RestartPolicy: dockerclient.RestartPolicy{
// 		//          Name:              "always",
// 		//          MaximumRetryCount: 0,
// 		//      },
// 		//  },
// 	}

// 	// pull images
// 	log.Info("Pulling image: %s\n", DockerImage)
// 	if err := docker.PullImage(DockerImage, nil); err != nil {
// 		log.Error("error pulling %s image: %s\n", DockerImage, err)
// 		return false, err
// 	}

// 	// start the show
// 	builderId, err := docker.CreateContainer(builderConfig, config.Repository)
// 	if err != nil {
// 		log.Error("error creating %s: %s\n", DockerImage, err)
// 		return false, err
// 	}
// 	var devices []dockerclient.DeviceMapping

// 	for i, c := range volumes {
// 		log.Info("Host path: " + c)
// 		log.Info("Container path:" + volumes[i])
// 		devices = append(devices, dockerclient.DeviceMapping{
// 			PathOnHost:      c,
// 			PathInContainer: volumes[i],
// 		})
// 	}
// 	os.Exit(1)

// 	hostConfig := &dockerclient.HostConfig{Devices: devices}
// 	if err := docker.StartContainer(builderId, hostConfig); err != nil {
// 		log.Error("error starting %s: %s\n", DockerImage, err)
// 		return false, err
// 	}
// 	docker.StartMonitorEvents(eventCallback, nil)

// 	// Remove handling -> gofunc to wait && delete
// 	// Valid case.
// 	//go func() {

// 	/// MEGLIO SENZA GOFUNC? in teoria dovrei attendere tra un job e un altro
// 	select {
// 	case wr := <-docker.Wait(builderId):
// 		docker.RemoveContainer(builderId, true, false)
// 		if wr.ExitCode == int(0) {
// 			// success, call the callback to save to db last valid commit
// 			return true, err
// 		} else {
// 			return false, err
// 		}
// 	}
// }
