package plugins

import (
	"time"

	"github.com/dikeert/dwmon-go/types"
	"github.com/jasonlvhit/gocron"
)

// Clock plugin is responsible for redring current time in specified format.
type ClockPlugin struct {
	format   string
	interval int
}

func (p *ClockPlugin) Initialize(flags types.Flags) {
	const format = "Mon, 02 Jan 2006, 03:04 PM"

	flags.StringVar(&p.format, "clock-format", format, "TBD")
	flags.IntVar(&p.interval, "clock-interval", 15, "TBD")
}

func (p *ClockPlugin) Start(s *gocron.Scheduler, notify chan bool) types.Module {
	f := func() {
		notify <- true
	}

	s.Every(15).Seconds().Do(f)

	return p.clock
}

func (p *ClockPlugin) clock(_ ...string) string {
	return time.Now().Format(p.format)
}
