package main

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
)

func (s *State) SendClientVersionDiff(channelID discord.ChannelID, cvd *ClientVersionDiff) error {
	log.Printf("Sending client version diff: %s %s", cvd.Old.ClientVersionUpload, cvd.New.ClientVersionUpload)

	embed := EmbedClientVersionDiff(cvd)

	if _, err := s.SendEmbeds(channelID, *embed); err != nil {
		return err
	}

	return nil
}
