package main

import (
	"log"

	"github.com/diamondburned/arikawa/v3/discord"
)

func (s *State) SendVersionDiff(channelID discord.ChannelID, vd *VersionDiff) error {
	log.Printf("Sending client version diff: %s %s", vd.Old.Upload, vd.New.Upload)

	embed := vd.Embed()

	if _, err := s.SendEmbeds(channelID, *embed); err != nil {
		return err
	}

	return nil
}
