package boson

// Processor is the interface for preprocessors, postprocessors and provisioners
type Processor interface {
	Process(*Build) ([]string, []string) // processor gets the workdir and the config file
	OnStart()
}
