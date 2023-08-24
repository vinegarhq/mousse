package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/vinegarhq/vinegar/util"
)

const VersionCheckURL = "https://clientsettingscdn.roblox.com/v2/client-version"

type Version struct {
	Version string `json:"version"`
	Upload  string `json:"clientVersionUpload"`
}

type (
	ChannelsVersions         map[string]Version
	BinariesChannelsVersions map[string]ChannelsVersions
)

func LatestVersion(binary string, channel string) (Version, error) {
	var ver Version
	url := VersionCheckURL + "/" + binary + "/channel/" + channel

	log.Printf("Fetching %s for channel %s (%s)", binary, channel, url)

	resp, err := util.Body(url)
	if err != nil {
		return Version{}, fmt.Errorf("failed to fetch version: %w", err)
	}

	err = json.Unmarshal([]byte(resp), &ver)
	if err != nil {
		return Version{}, fmt.Errorf("failed to unmarshal clientsettings: %w", err)
	}

	if ver.Upload == "" {
		return Version{}, errors.New("no version found")
	}

	return ver, nil
}

func ChannelsLatestVersions(binary string) (ChannelsVersions, error) {
	cvs := make(ChannelsVersions, 0)

	log.Printf("Fetching all latest versions for binary %s", binary)

	for _, c := range Channels {
		v, err := LatestVersion(binary, c)
		if err != nil {
			return ChannelsVersions{}, err
		}

		cvs[c] = v
	}

	return cvs, nil
}

func BinariesChannelsLatestVersions() (BinariesChannelsVersions, error) {
	bcvs := make(BinariesChannelsVersions, 0)

	log.Println("Fetching all latest versions for all binaries")

	for _, b := range Binaries {
		bcv, err := ChannelsLatestVersions(b)
		if err != nil {
			return BinariesChannelsVersions{}, err
		}

		bcvs[b] = bcv
	}

	return bcvs, nil
}
