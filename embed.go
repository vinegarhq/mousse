package main

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
)

func (vd *VersionDiff) Embed() *discord.Embed {
	embed := discord.NewEmbed()

	embed.Title = fmt.Sprintf("%s@%s", vd.Binary, vd.Channel)
	embed.Description = fmt.Sprintf(
		"```diff\n- %s (%s)\n+ %s (%s)\n```\n",
		vd.Old.Upload, vd.Old.Version,
		vd.New.Upload, vd.New.Version,
	)
	embed.Color = 0xAFC147

	return embed
}
