package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/dikeert/dwmon-go/internal/plugins"
	"github.com/dikeert/dwmon-go/internal/sinks"
	"github.com/dikeert/dwmon-go/types"
	"github.com/jasonlvhit/gocron"
	"github.com/spf13/cobra"
)

var oobPlugins = map[string]types.Plugin{
	"clock":  &plugins.ClockPlugin{},
	"echo":   &plugins.EchoPlugin{},
	"mpd":    &plugins.MpdPlugin{},
	"wakeup": &plugins.WakeupPlugin{},
	"shell":  &plugins.ShellPlugin{},
}

var oobSinks = map[string]types.Sink{
	"stdout":   sinks.Stdout,
	"xsetroot": sinks.Xsetroot,
}

type output struct {
	strings.Builder
	sink types.Sink
}

func (o *output) dump() error {
	err := o.sink(o.String())
	if err == nil {
		o.Reset()
	}

	return err
}

type command struct {
	cobra.Command
	scheduler *gocron.Scheduler
	updates   chan bool
	enabled   *[]string // list of enabled plugins, populated from cmdline
	format    *string
	sink      *string
}

func (c *command) run(cmd *cobra.Command, args []string) error {
	modules, err := start(c, c.scheduler, c.updates)
	if err != nil {
		return err
	}

	tmpl, err := template.New("output").Funcs(modules).Parse(*c.format)
	if err != nil {
		return err
	}

	sink, err := getSink(c)
	if err != nil {
		return err
	}

	go func() {
		for range c.updates {
			var err error
			if err = tmpl.Execute(sink, []string{}); err != nil {
				fmt.Fprintf(os.Stderr, "Error while rendering: %s", err)
				os.Exit(1)
			}

			if err = sink.dump(); err != nil {
				fmt.Fprintf(os.Stderr, "Error while rendering: %s", err)
				os.Exit(1)
			}

		}
	}()

	c.updates <- true // first render

	<-c.scheduler.Start()

	return nil
}

func main() {
	cmd := createCommand()

	updates := make(chan bool, 10)
	s := gocron.NewScheduler()

	cmd.scheduler = s
	cmd.updates = updates
	cmd.RunE = cmd.run

	cmd.SetArgs(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func getSink(cmd *command) (*output, error) {
	sinkName := *cmd.sink

	if sink, ok := oobSinks[sinkName]; ok {
		return &output{sink: sink}, nil
	}

	return nil, fmt.Errorf("Sink '%s' does not exist", sinkName)
}

func createCommand() *command {
	cmd := &command{
		Command: cobra.Command{
			Use:   "dwmon",
			Short: "Output status information",
			Long:  "TBD",
			Args:  cobra.ArbitraryArgs,
		},
	}

	// TODO: usage (in TBD)
	cmd.enabled = cmd.Flags().StringSliceP("plugins", "p", []string{}, "TBD")

	// TODO: usage (in TBD)
	cmd.format = cmd.Flags().StringP("format", "f", "", "TBD")

	// TODO: usage (in TBD)
	cmd.sink = cmd.Flags().StringP("sink", "s", "stdout", "TBD")

	initialize(cmd)

	return cmd
}

func initialize(cmd *command) {
	for _, plugin := range oobPlugins {
		plugin.Initialize(cmd.Flags())
	}
}

func start(cmd *command, s *gocron.Scheduler, updates chan bool) (map[string]interface{}, error) {
	modules := make(map[string]interface{})

	enabled := *cmd.enabled

	if len(enabled) < 1 {
		return nil, errors.New("No plugins are enabled")
	}

	for _, name := range enabled {
		if plugin, ok := oobPlugins[name]; ok {
			modules[name] = plugin.Start(s, updates)
		} else {
			return nil, fmt.Errorf("There is no plugin '%s'", name)
		}
	}

	return modules, nil
}
