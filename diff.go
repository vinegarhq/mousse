package main

import (
	"log"
)

type VersionDiff struct {
	Channel string
	Binary  string
	Old     *Version
	New     *Version
}

func (bcvs *BinariesChannelsVersions) Check(onMismatch func(*VersionDiff) error) {
	log.Printf("Checking for updates of all binaries")

	for _, b := range Binaries {
		bcv := (*bcvs)[b]
		bcv.Check(b, onMismatch)
	}
}

func (cvs *ChannelsVersions) Check(binary string, onMismatch func(*VersionDiff) error) {
	log.Printf("Checking for updates of %s", binary)

	for _, c := range Channels {
		cv := (*cvs)[c]
		ver, err := LatestVersion(binary, c)
		if err != nil {
			log.Println(err)

			continue
		}

		if ver.Upload == cv.Upload {
			continue
		}

		err = onMismatch(&VersionDiff{
			Channel: c,
			Binary:  binary,
			New:     &ver,
			Old:     &cv,
		})

		if err != nil {
			log.Println(err)

			continue
		}

		(*cvs)[c] = ver
	}
}
