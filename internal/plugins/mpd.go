package plugins

import (
	"fmt"
	"os"

	"github.com/dikeert/dwmon-go/types"
	"github.com/jasonlvhit/gocron"
	"github.com/vincent-petithory/mpdclient"
)

func (p *MpdPlugin) mpd(params ...string) string {
	if len(params) < 1 {
		params = []string{"Title"}
	}

	output := ""

	for _, param := range params {
		var token string

		if val, ok := p.status[param]; ok {
			token = val
		} else {
			token = param
		}

		output = fmt.Sprintf("%s%s", output, token)
	}

	if p.maxLen > 0 {
		return output[0:p.maxLen]
	}

	return output
}

// The only purpose of Wake Up plugin is to cause re-rendering of the status
// line whenever user demands it by sending USR1 to the process.
type MpdPlugin struct {
	host   string
	port   uint
	maxLen uint
	client *mpdclient.MPDClient
	status map[string]string
}

func (p *MpdPlugin) Initialize(f types.Flags) {
	//TODO: proper doc
	f.StringVar(&p.host, "mpd-host", "localhost", "TBD")
	//TODO: proper doc
	f.UintVar(&p.port, "mpd-port", 6600, "TBD")
	f.UintVar(&p.maxLen, "mpd-max-length", 0, "TBD")
	p.status = make(map[string]string)
}

func (p *MpdPlugin) Start(_ *gocron.Scheduler, updates chan bool) types.Module {
	client, err := connect(p.host, p.port)
	if err == nil {
		events := idle(client)
		p.client = client
		go listen(p.client, &p.status, events, updates)
	} else {
		p.client = nil
		fmt.Fprintf(os.Stderr, "Unable to connect to mpd: %s", err)
	}

	return p.mpd
}

func connect(host string, port uint) (*mpdclient.MPDClient, error) {
	mpdc, err := mpdclient.Connect(host, port)
	if err != nil {
		return nil, err
	}

	return mpdc, nil
}

func idle(client *mpdclient.MPDClient) chan string {
	events := client.Idle("player")
	return events.Ch
}

func listen(client *mpdclient.MPDClient, status *map[string]string,
	events chan string, updates chan bool) {
	for subsystem := range events {
		switch subsystem {
		case "player":
			info, err := client.CurrentSong()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to fetch current song: %s", err)
				continue
			}

			for _, field := range []string{"Title", "Artist", "Album", "AlbumArtist"} {
				(*status)[field] = (*info)[field]
			}

			updates <- true
		}
	}
}
