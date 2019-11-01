package plugins

import (
	"os/exec"
	"strings"

	"github.com/dikeert/mon-go/types"
	"github.com/jasonlvhit/gocron"
)

func shell(args ...string) string {
	output, err := exec.Command(args[0], args[1:]...).Output()
	if err == nil {
		return strings.Trim(string(output), "\r\n")
	}

	return err.Error()
}

type ShellPlugin struct{}

func (p *ShellPlugin) Initialize(_ types.Flags) {}

func (p *ShellPlugin) Start(_ *gocron.Scheduler, _ chan bool) types.Module {
	return shell
}
