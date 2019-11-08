package plugins

import (
	"fmt"
	"math"
	"os"

	"github.com/dikeert/dwmon-go/types"
	"github.com/jasonlvhit/gocron"
	"github.com/vincent-petithory/mpdclient"
)

var properties = []string{
	"Title",
	"Artist",
	"Album",
	"AlbumArtist",
}

var subsystems = map[string]interface{}{
	"player": nil,
}

func (p *MpdPlugin) mpd(params ...string) string {
	if len(params) < 1 {
		params = properties[0:1]
	}

	output := ""

	for _, param := range params {
		var property string

		if val, ok := p.status[param]; ok {
			property = val
		} else {
			property = param
		}

		output = fmt.Sprintf("%s%s", output, property)
	}

	if p.maxLen > 0 {
		length := math.Min(float64(p.maxLen), float64(len(output)))
		return output[0:int64(length)]
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
		p.client = client
		p.status = status(client)

		go listen(p, updates)
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

func status(c *mpdclient.MPDClient) map[string]string {
	info, err := c.CurrentSong()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch current song: %s", err)
		return map[string]string{}
	}

	return *info
}

func listen(p *MpdPlugin, updates chan bool) {
	for subsystem := range idle(p.client) {
		if _, ok := subsystems[subsystem]; ok {
			p.status = status(p.client)
			updates <- true
		}
	}
}

func idle(client *mpdclient.MPDClient) chan string {
	wanted := []string{}
	for key := range subsystems {
		wanted = append(wanted, key)
	}

	return client.Idle(wanted...).Ch
}
