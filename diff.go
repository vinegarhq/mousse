package main

import (
	"fmt"
	"log"

	"github.com/vinegarhq/vinegar/roblox"
	"github.com/vinegarhq/vinegar/roblox/api"
)

var Binaries = []roblox.BinaryType{
	roblox.Player,
	roblox.Studio,
}

var Channels = []string{
	"LIVE",
}

type Version struct {
	GUID  string
	Canon string
}

type (
	ChannelsVersions         map[string]Version
	BinariesChannelsVersions map[roblox.BinaryType]ChannelsVersions
)

type VersionDiff struct {
	Binary  roblox.BinaryType
	Channel string
	Old     *Version
	New     *Version
}

func (vd VersionDiff) String() string {
	return fmt.Sprintf("%s: %s, %s -> %s", vd.Binary, vd.Channel, vd.Old.GUID, vd.New.GUID)
}

type VersionDiffMismatchFunc func(*VersionDiff) error

func NewBinariesChannelsVersions() BinariesChannelsVersions {
	bcvs := make(BinariesChannelsVersions)

	for _, b := range Binaries {
		bcvs[b] = make(ChannelsVersions)
	}

	return bcvs
}

func (bcvs BinariesChannelsVersions) Check(fn VersionDiffMismatchFunc) {
	for _, b := range Binaries {
		bcvs[b].Check(b, fn)
	}
}

func (cvs ChannelsVersions) Check(bt roblox.BinaryType, fn VersionDiffMismatchFunc) {
	log.Printf("Checking for version changes for %s for all channels", bt)

	for _, c := range Channels {
		cv := cvs[c]

		cvu, err := api.GetClientVersion(bt, c)
		if err != nil {
			log.Printf("%s: channel %s: %s", bt, c, err)

			continue
		}

		if cvu.ClientVersionUpload == cvs[c].GUID {
			continue
		}

		nv := Version{
			GUID:  cvu.ClientVersionUpload,
			Canon: cvu.Version,
		}

		err = fn(&VersionDiff{
			Binary:  bt,
			Channel: c,
			New:     &nv,
			Old:     &cv,
		})

		cvs[c] = nv

		if err != nil {
			log.Println(err)

			continue
		}
	}
}
