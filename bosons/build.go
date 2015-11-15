// Package boson includes the boson's builds, builders, processors utilities
package boson

import "github.com/mudler/boson/shared/utils"

// Build the build type have a *utils.Config, a Commit  and a PrevCommit hash
type Build struct {
	Config     *utils.Config
	Commit     string
	PrevCommit string
	Extras     []string
}
