package plugins

import (
	"strings"

	"github.com/dikeert/mon-go/types"
	"github.com/jasonlvhit/gocron"
)

func echo(args ...string) string {
	return strings.Join(args, ",")
}

// Echo Plugin recieves and input and prints it into
// status output as is
type EchoPlugin struct{}

func (p *EchoPlugin) Initialize(_ types.Flags) {}

func (p *EchoPlugin) Start(_ *gocron.Scheduler, _ chan bool) types.Module {
	return echo
}
