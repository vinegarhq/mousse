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
	ccvs, err := AllLatestVersions()
	if err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(2 * time.Minute)

		for _, c := range Channels {
			ccv := ccvs[c]
			cv, err := LatestVersion(c)
			if err != nil {
				log.Println(err)

				continue
			}

			if cv.ClientVersionUpload == ccv.ClientVersionUpload {
				continue
			}

			log.Println("Mismatch (%s): %s %s", c, cv.ClientVersionUpload, ccv.ClientVersionUpload)

			err = s.SendClientVersionDiff(discord.ChannelID(channelID), &ClientVersionDiff{
				Channel: c,
				New:     &cv,
				Old:     &ccv,
			})
			if err != nil {
				log.Println(err)
				continue
			}

			ccvs[c] = cv
		}
	}
}
