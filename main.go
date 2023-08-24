package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

var (
	token     string
	channelID int64
)

var Binaries = []string{
	"WindowsPlayer",
	"WindowsStudio",
	"WindowsStudio64",
	"MacPlayer",
	"MacStudio",
}

var Channels = []string{
	"LIVE",
	"ZIntegration",
	"ZCanary",
	"ZFlag",
	"ZNext",
}

type State struct {
	*state.State
}

func init() {
	flag.StringVar(&token, "token", "", "Discord Bot Token")
	flag.Int64Var(&channelID, "channel", 0, "Channel ID")
}

func main() {
	flag.Parse()
	log.Println("Starting RDCW")

	s := &State{
		State: state.New("Bot " + token),
	}

	s.AddIntents(gateway.IntentGuilds)

	if err := s.Open(context.TODO()); err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	// first run
	bcvs, err := BinariesChannelsLatestVersions()
	if err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(2 * time.Minute)

		bcvs.Check(func(vd *VersionDiff) error {
			return s.SendVersionDiff(discord.ChannelID(channelID), vd)
		})
	}
}
