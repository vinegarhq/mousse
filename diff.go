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

func (cvs *ChannelsVersions) Check(binary string, onMismatch func(*VersionDiff) error) {
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
