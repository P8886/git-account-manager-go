//go:build !windows

package gitops

import "os/exec"

func prepareCmd(cmd *exec.Cmd) {
	// No special handling needed for Unix-like systems
}
