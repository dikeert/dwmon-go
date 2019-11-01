package types

import (
	"github.com/jasonlvhit/gocron"
)

type Module func(...string) string
type Sink func(string) error

type Flags interface {
	StringVar(*string, string, string, string)
	IntVar(*int, string, int, string)
}

type Plugin interface {
	Initialize(Flags)
	Start(*gocron.Scheduler, chan bool) Module
}
