// Plugin register themselves here, the registry keep tracks of plugins to redirect the messages
package plugin_registry

import (
	"github.com/mudler/boson/processor"

	"reflect"
	"strings"
)

// These are registered plugins
var Preprocessors = map[string]processor.Processor{}
var Postprocessors = map[string]processor.Processor{}

// Register a Preprocessor
func RegisterPreprocessor(p processor.Processor) {
	Preprocessors[KeyOf(p)] = p
}

// Register a Postprocessor
func RegisterPostprocessor(p processor.Processor) {
	Postprocessors[KeyOf(p)] = p
}

func KeyOf(p processor.Processor) string {
	return strings.TrimPrefix(reflect.TypeOf(p).String(), "*")
}
