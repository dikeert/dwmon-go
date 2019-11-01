package plugins

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dikeert/mon-go/types"
	"github.com/jasonlvhit/gocron"
)

func wakeup(_ ...string) string {
	return ""
}

// The only purpose of Wake Up plugin is to cause re-rendering of the status
// line whenever user demands it by sending USR1 to the process.
type WakeupPlugin struct{}

func (p *WakeupPlugin) Initialize(f types.Flags) {
}

func (p *WakeupPlugin) Start(_ *gocron.Scheduler, updates chan bool) types.Module {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGUSR1)

	go func() {
		for sig := range sigc {
			switch sig {
			case syscall.SIGUSR1:
				updates <- true
			}
		}
	}()

	return wakeup
}
