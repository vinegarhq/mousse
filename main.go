package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	token     string
	channelID string
)

var Binaries = []BinaryType{
	WindowsPlayer,
	WindowsStudio,
	WindowsStudio64,
	MacPlayer,
	MacStudio,
}

var Channels = []string{
	"LIVE",
	"ZIntegration",
	"ZCanary",
	"ZFlag",
	"ZNext",
}

func init() {
	flag.StringVar(&token, "token", "", "Discord Bot Token")
	flag.StringVar(&channelID, "channel", "1143583777831010394", "Channel ID")
}

func main() {
	flag.Parse()
	log.Println("Starting Mousse")

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	dg.LogLevel = 1
	dg.Identify.Intents = discordgo.IntentsGuilds

	if err := dg.Open(); err != nil {
		log.Fatal(err)
	}

	log.Println("Mousse is now running. Send TERM/INT to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-sc
		dg.Close()
		os.Exit(0)
	}()

	err = dg.UpdateWatchStatus(0, "Roblox's binaries")
	if err != nil {
		dg.Close()
		log.Fatal(err)
	}

	// first run
	bcvs := make(BinariesChannelsVersions, 0)
	bcvs.Check(func(vd *VersionDiff) error {
		return nil
	})

	for {
		time.Sleep(2 * time.Minute)

		bcvs.Check(func(vd *VersionDiff) error {
			log.Printf("Sending version embed diff: %s %s", vd.Old.GUID, vd.New.GUID)

			if _, err := dg.ChannelMessageSendEmbed(channelID, vd.Embed()); err != nil {
				return err
			}

			return nil
		})
	}
}

func (vd *VersionDiff) Embed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s@%s", vd.Binary, vd.Channel),
		Description: fmt.Sprintf(
			"```diff\n- %s (%s)\n+ %s (%s)\n```\n",
			vd.Old.Real, vd.Old.GUID,
			vd.New.Real, vd.New.GUID,
		),
		Color: 0xAFC147,
	}
}
