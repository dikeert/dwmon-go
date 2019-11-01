package sinks

import "os/exec"

func Xsetroot(output string) error {
	cmd := exec.Command("xsetroot", "-name", output)
	return cmd.Run()
}
