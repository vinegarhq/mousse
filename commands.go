package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/robloxapi/rbxdhist"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
)

// Taken from github.com/diamondburned/arikawa/blob/v3/0-examples

var commands = []api.CreateCommandData{
	{
		Name:        "status",
		Description: "Display current tracking channels and binaries",
	},
	{
		Name:        "expired",
		Description: "Check if a given deployment GUID is expired",
		Options: discord.CommandOptions{
			&discord.StringOption{
				OptionName:  "guid",
				Description: "Deployment GUID to check if expired",
				Required:    true,
			},
		},
	},
	{
		Name:        "query",
		Description: "Retrieve general information about a GUID",
		Options: discord.CommandOptions{
			&discord.StringOption{
				OptionName:  "guid",
				Description: "Deployment GUID to query",
				Required:    true,
			},
		},
	},
}

func cmdStatus(_ context.Context, _ cmdroute.CommandData) *api.InteractionResponseData {
	return &api.InteractionResponseData{Content: option.NewNullableString(
		fmt.Sprintf("Tracking binaries `%s` with channels `%s`", Binaries, Channels),
	)}
}

func cmdQuery(_ context.Context, d cmdroute.CommandData) *api.InteractionResponseData {
	var o struct {
		GUID string `discord:"guid"`
	}

	if err := d.Options.Unmarshal(&o); err != nil {
		return errorResponse(err)
	}

	j, err := findJob(o.GUID)
	if err != nil {
		return errorResponse(err)
	}

	e := discord.NewEmbed()
	e.Title = j.GUID
	e.Fields = []discord.EmbedField{
		{Name: "Build", Value: j.Build},
		{Name: "Time", Value: j.Time.String()},
		{Name: "Version", Value: fmt.Sprintf("`%s`", j.Version)},
		{Name: "GitHash", Value: j.GitHash},
	}

	return &api.InteractionResponseData{
		Embeds: &[]discord.Embed{
			*e,
		},
	}
}

func errorResponse(err error) *api.InteractionResponseData {
	return &api.InteractionResponseData{
		Content:         option.NewNullableString("**Error:** " + err.Error()),
		Flags:           discord.EphemeralMessage,
		AllowedMentions: &api.AllowedMentions{ /* none */ },
	}
}

func findJob(guid string) (*rbxdhist.Job, error) {
	r, err := http.Get("https://setup.rbxcdn.com/DeployHistory.txt")
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, err
	}

	t := rbxdhist.Lex(b)
	for _, s := range t {
		j, ok := s.(*rbxdhist.Job)
		if !ok {
			continue
		}

		if j.GUID != guid {
			continue
		}

		return j, nil
	}

	return nil, errors.New("GUID not found in deployment history")
}
