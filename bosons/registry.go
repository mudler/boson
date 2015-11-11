// Package boson register plugins here, the registry keep tracks of plugins to redirect the messages
package boson

import (
	"reflect"
	"strings"
)

// Preprocessors contains a map of Preprocessor
var Preprocessors = map[string]Processor{}

// Postprocessors contains a map of Postprocessors
var Postprocessors = map[string]Processor{}

// Provisioners contains a map of Provisioners
var Provisioners = map[string]Processor{}

// RegisterPreprocessor Registers a Preprocessor
func RegisterPreprocessor(p Processor) {
	Preprocessors[keyOf(p)] = p
}

// RegisterPostprocessor Registers a Postprocessor
func RegisterPostprocessor(p Processor) {
	Postprocessors[keyOf(p)] = p
}

// RegisterProvisioner Registers a Provisioner
func RegisterProvisioner(p Processor) {
	Provisioners[keyOf(p)] = p
}

func keyOf(p Processor) string {
	return strings.TrimPrefix(reflect.TypeOf(p).String(), "*")
}
