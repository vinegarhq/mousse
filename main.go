package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

var (
	token     string
	channelID int64
)

func init() {
	flag.StringVar(&token, "token", "", "Discord Bot Token")
	flag.Int64Var(&channelID, "channel", 0, "Channel ID")
}

func main() {
	flag.Parse()
	log.Println("Starting Mousse")

	r := cmdroute.NewRouter()
	r.AddFunc("status", cmdStatus)
	r.AddFunc("query", cmdQuery)

	s := state.New("Bot " + token)
	s.AddInteractionHandler(r)
	s.AddIntents(gateway.IntentGuilds)

	if err := cmdroute.OverwriteCommands(s, commands); err != nil {
		log.Fatalln("cannot update commands:", err)
	}

	if err := s.Open(context.Background()); err != nil {
		log.Fatalln("cannot open:", err)
	}

	err := s.Gateway().Send(context.TODO(), &gateway.UpdatePresenceCommand{
		Activities: []discord.Activity{{
			Name: "Roblox's binaries",
			Type: discord.WatchingActivity,
		}},
	})
	if err != nil {
		log.Println("cannot update activity:", err)
	}

	log.Println("Mousse is now running. Send TERM/INT to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-sc
		s.Close()
		os.Exit(0)
	}()

	bcvs := NewBinariesChannelsVersions()
	bcvs.Check(func(vd *VersionDiff) error {
		log.Printf("First run %s: %s, %s", vd.Binary, vd.Channel, vd.New.GUID)
		return nil
	})

	for {
		time.Sleep(2 * time.Minute)

		bcvs.Check(func(vd *VersionDiff) error {
			log.Printf("Sending version embed diff: %s", vd)

			if _, err := s.SendEmbeds(discord.ChannelID(channelID), *vd.Embed()); err != nil {
				return err
			}

			return nil
		})
	}
}

func (vd *VersionDiff) Embed() *discord.Embed {
	embed := discord.NewEmbed()

	embed.Title = fmt.Sprintf("%s@%s", vd.Binary, vd.Channel)
	embed.Description = fmt.Sprintf(
		"```diff\n- %s (%s)\n+ %s (%s)\n```\n",
		vd.Old.Canon, vd.Old.GUID,
		vd.New.Canon, vd.New.GUID,
	)

	embed.Color = 0xAFC147
	if vd.Channel == "LIVE" {
		embed.Color = 0xCC241D
	}

	return embed
}
