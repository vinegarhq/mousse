package main

import (
	"log"
)

type (
	ChannelsVersions         map[string]Version
	BinariesChannelsVersions map[BinaryType]ChannelsVersions
)

type VersionDiff struct {
	Channel string
	Binary  BinaryType
	Old     *Version
	New     *Version
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

func (cvs ChannelsVersions) Check(bt BinaryType, fn VersionDiffMismatchFunc) {
	log.Printf("Checking for version changes for %s for all channels", bt)

	for _, c := range Channels {
		cv := cvs[c]

		ver, err := bt.Version(c)
		if err != nil {
			log.Printf("%s: channel %s: %s", bt, c, err)

			continue
		}

		if ver.GUID == cv.GUID {
			continue
		}

		err = fn(&VersionDiff{
			Channel: c,
			Binary:  bt,
			New:     &ver,
			Old:     &cv,
		})

		if err != nil {
			log.Println(err)

			continue
		}

		cvs[c] = ver
	}
}
