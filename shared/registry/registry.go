// Package plugin_registry register plugins here, the registry keep tracks of plugins to redirect the messages
package plugin_registry

import (
	"github.com/mudler/boson/processor"

	"reflect"
	"strings"
)

// Preprocessors contains a map of Preprocessor
var Preprocessors = map[string]processor.Processor{}

// Postprocessors contains a map of Postprocessors
var Postprocessors = map[string]processor.Processor{}

// Provisioners contains a map of Provisioners
var Provisioners = map[string]processor.Processor{}

// RegisterPreprocessor Registers a Preprocessor
func RegisterPreprocessor(p processor.Processor) {
	Preprocessors[keyOf(p)] = p
}

// RegisterPostprocessor Registers a Postprocessor
func RegisterPostprocessor(p processor.Processor) {
	Postprocessors[keyOf(p)] = p
}

// RegisterProvisioner Registers a Provisioner
func RegisterProvisioner(p processor.Processor) {
	Provisioners[keyOf(p)] = p
}

func keyOf(p processor.Processor) string {
	return strings.TrimPrefix(reflect.TypeOf(p).String(), "*")
}
