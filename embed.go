package main

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
)

func EmbedClientVersionDiff(cvd *ClientVersionDiff) *discord.Embed {
	embed := discord.NewEmbed()

	embed.Title = fmt.Sprintf("WindowsPlayer@%s", cvd.Channel.String())
	embed.Description = fmt.Sprintf(
		"```diff\n- %s (%s)\n+ %s (%s)\n```\n",
		cvd.Old.ClientVersionUpload, cvd.Old.Version,
		cvd.New.ClientVersionUpload, cvd.New.Version,
	)
	embed.Color = 0xAFC147

	return embed
}
