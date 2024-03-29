package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/apprehensions/rbxbin"
	cs "github.com/apprehensions/rbxweb/clientsettings"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

func main() {
	t := flag.String("token", "", "Discord Bot Token")
	c := flag.Int64("channel", 1143583777831010394, "Channel ID")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	errorResponse := func(err error) *api.InteractionResponseData {
		return &api.InteractionResponseData{
			Content:         option.NewNullableString("**Error:** " + err.Error()),
			Flags:           discord.EphemeralMessage,
			AllowedMentions: &api.AllowedMentions{ /* none */ },
		}
	}

	r := cmdroute.NewRouter()
	r.AddFunc("query", func(_ context.Context, d cmdroute.CommandData) *api.InteractionResponseData {
		var o struct {
			GUID string `discord:"guid"`
		}

		if err := d.Options.Unmarshal(&o); err != nil {
			return errorResponse(err)
		}

		m, err := rbxbin.GetMirror()
		if err != nil {
			return errorResponse(err)
		}

		js, err := m.Jobs()
		if err != nil {
			return errorResponse(err)
		}

		var f *rbxbin.Job
		for _, j := range js {
			if j.GUID == o.GUID {
				f = j
			}
		}
		if f == nil {
			return errorResponse(errors.New("GUID not found in deployment history"))
		}

		e := discord.NewEmbed()
		e.Title = f.GUID
		e.Fields = []discord.EmbedField{
			{Name: "Build", Value: f.Build},
			{Name: "Time", Value: f.Time.String()},
			{Name: "Version", Value: fmt.Sprintf("`%s`", f.Version)},
			{Name: "GitHash", Value: f.GitHash},
		}

		return &api.InteractionResponseData{Embeds: &[]discord.Embed{*e}}
	})

	s := state.New("Bot " + *t)
	s.AddInteractionHandler(r)
	s.AddIntents(gateway.IntentGuilds)

	if err := cmdroute.OverwriteCommands(s, []api.CreateCommandData{{
		Name:        "query",
		Description: "Retrieve general information about a GUID",
		Options: discord.CommandOptions{&discord.StringOption{
			OptionName:  "guid",
			Description: "Deployment GUID to query",
			Required:    true,
		}},
	}}); err != nil {
		log.Fatalln("cannot update commands:", err)
	}

	s.AddHandler(func(*gateway.ReadyEvent) {
		if err := s.Gateway().Send(ctx, &gateway.UpdatePresenceCommand{
			Activities: []discord.Activity{{
				Name: "Roblox's binaries",
				Type: discord.WatchingActivity,
			}},
		}); err != nil {
			log.Fatalln("cannot update activity:", err)
		}
		slog.Info("Connected!")
	})

	k := make(map[cs.BinaryType]*cs.ClientVersion)
	for _, b := range []cs.BinaryType{cs.WindowsPlayer, cs.WindowsStudio64} {
		k[b] = nil
	}

	go func() {
		for {
			for b, v := range k {
				cv, err := cs.GetClientVersion(b, "")
				if err != nil {
					slog.Error("Failed to fetch ClientVersion", "error", err)
					continue
				}
				k[b] = cv

				slog.Info("Fetched ClientVersion", "binary", b, "version", cv)

				if v == nil || cv.GUID == v.GUID {
					continue
				}

				e := discord.NewEmbed()
				e.Color = 0xCC241D
				e.Title = fmt.Sprintf("%s@%s", b.Short(), "LIVE")
				e.Description = fmt.Sprintf(
					"```diff\n- %s (%s)\n+ %s (%s)\n```\n",
					v.Version, v.GUID, cv.Version, cv.GUID,
				)

				slog.Info("Sending updated ClientVersion", "binary", b, "old_version", v, "new_version", cv)
				if _, err := s.SendEmbeds(discord.ChannelID(*c), *e); err != nil {
					slog.Error("Failed to send update", "error", err)
					continue
				}
			}

			time.Sleep(2 * time.Minute)
		}
	}()

	if err := s.Connect(ctx); err != nil {
		log.Fatalln("connection closed:", err)
	}
}
