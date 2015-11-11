# Boson

[![wercker status](https://app.wercker.com/status/5d8c39eb339bbff485baab3840a4a099/m "wercker status")](https://app.wercker.com/project/bykey/5d8c39eb339bbff485baab3840a4a099)

### Polling CI / Provision system for buildchains

Boson is a CI System that uses POLL technique instead of WebHooks. Yes, a question now will pop in your mind

*Why polling?*

Because there are situations where you just can't setup a webhook.

*Why don't fork and track the repository then?*

Mantaining can be expensive, or even worst, you are like me. **Lazy**.


### Reinventing again the wheel (?),  but with simplicity in mind:


First time, Boson will sync up with your repo, so let's say that your last commit is this: *bc97686f0c710bf19a0f391d2cb49c7237814b73*, if it was the first start it will take this commit as base point

    Boson -> bc97686f0c710bf19a0f391d2cb49c7237814b73
    Repository -> bc97686f0c710bf19a0f391d2cb49c7237814b73

If someone commits *1fd4f077c760139886afc0da52674faebcb8a926* it will be like this:

    Boson -> bc97686f0c710bf19a0f391d2cb49c7237814b73
    Repository -> bc97686f0c710bf19a0f391d2cb49c7237814b73 -> 1fd4f077c760139886afc0da52674faebcb8a926

now depending on your specified PreProcessor, Boson will perform a custom action between *bc97686f0c710bf19a0f391d2cb49c7237814b73* and *1fd4f077c760139886afc0da52674faebcb8a926* commits.
For example, the Gentoo PreProcessor will check for touched ebuilds and give as result a list of packages.
The PreProcessor then construct and orchestrate the arguments and volumes to feed on the other phase: the spawning of a throw-away docker container.

In such way, the Gentoo Preprocessor can spawn a builder machine that compile and tests your files.

### Usage

Compile it as a normal go project or [download the precompiled binary in the release page](https://github.com/mudler/boson/releases) and then launch it

	./boson -c your-boson-file.yaml

### Boson file


Boson file lets you define few options for now

    ---
    repository: https://github.com/yourusername/yourrepo.git
    docker_image: dockerimagename
    preprocessor: gentoo.Gentoo # as for now is the only available
    postprocessor:
        - Some
        - Some
    polltime: 50 # in seconds
    artifacts_dir: /somewhere/yourstash/artifacts # preprocessor can address this to their formal outputù
    volumes:
        - somedirinhost:somedirincontainer:ro
    env:
        - SOMETHING=somethingelse
    args:
        - --verbose
    log_dir: /some/logpath/dir
    separate_artifacts: true # this will generate for each commit a new directory with your container output specified in 'artifacts_dir'
    tmpdir: /tmp/boson #your tempfile, defaults is /var/tmp/boson


* **repository**: Contains the git repository to clone
* **docker_image**: The docker image name (that can be pulled from DockerHub)
* **preprocessor**: Plugin called to evaluate the repository difference
* **polltime**: Seconds to pass between one pull time and another (the builds, as for now is blocking. This is because after a job, it will be evaluated the next "portion") - in case , for example when the machine is building will be committed one or more commit, the next round it will be evaluated since the last commit built
* **artifacts_dir**: Typically this process involve an output dir (i'll discuss this later on my usage scenario)
* **volumes**: A list of volumes to be mounted o the container
* **env**: Alist of env to be set in the container
* **args**: Optional args to pass as argument at the container when running it
* **log_dir**: Where to store your container output after build
* **separate_artifacts**: If this enabled will be used "artifacts_dir" as parent of directory tree containing your job output for each commit
* **tmpdir**: Specify a different temporary directory

Boson with providers and commit

....

#### Example (my personal use case)

##### Use case: Setting up a CI Buildbox for the packages pushed in the dev's overlay:

    ---
    repository: https://github.com/Sabayon/for-gentoo.git
    docker_image: sabayon/builder-amd64
    preprocessor: gentoo.Gentoo
    polltime: 50
    artifacts_dir: /mydata/for-gentoo/artifacts
    log_dir: /mydata/for-gentoo/logs
    separate_artifacts: true
    tmpdir: /tmp/boson

In such case the *sabayon/builder-amd64* entrypoint is a wrapper script to emerge that actually ensures and gives other options (to make the process simpler).

The **gentoo.Gentoo** *PreProcessor* will process the git portion of the build and extract from the diff the ebuilds that needs to be compiled (and it's versions) as a List given to the Docker container as argument.

The **artifacts_dir** is mounted automatically from **gentoo.Gentoo** to */usr/portage/packages* wich is the default emerge output for builded packages.

So at each commit the image take care to test *without* carrying the actual steps of the test in a yaml file (**for now**) but instead uses the container entrypoint. I still didn't implemented a yaml testing file definition, leaving the already-defined logic in the repos: This is to mean that if the repository already runs for example on drone.io, you can also call drone itself and leave the test logic to an another engine.

This allows us to track also external github repository where we don't have permission to set a webhook, and simply creating a custom docker image to implement our testing logic against it.



##### Use case: Setting up a Matter Buildbox that can commit automatically to overlays

...
