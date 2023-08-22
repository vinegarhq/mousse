package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/vinegarhq/vinegar/util"
)

const VersionCheckURL = "https://clientsettingscdn.roblox.com/v2/client-version"

type ClientVersion struct {
	Version                 string `json:"version"`
	ClientVersionUpload     string `json:"clientVersionUpload"`
	BootstrapperVersion     string `json:"bootstrapperVersion"`
	NextClientVersionUpload string `json:"nextClientVersionUpload,omitempty"`
	NextClientVersion       string `json:"nextClientVersion,omitempty"`
}

type ClientVersionDiff struct {
	Channel
	Old *ClientVersion
	New *ClientVersion
}

type ChannelsClientVersions map[Channel]ClientVersion

func LatestVersion(channel Channel) (ClientVersion, error) {
	var cv ClientVersion

	url := VersionCheckURL + "/WindowsPlayer/channel/" + channel.String()

	log.Println(url)

	resp, err := util.Body(url)
	if err != nil {
		return ClientVersion{}, fmt.Errorf("failed to fetch version: %w", err)
	}

	err = json.Unmarshal([]byte(resp), &cv)
	if err != nil {
		return ClientVersion{}, fmt.Errorf("failed to unmarshal clientsettings: %w", err)
	}

	if cv.ClientVersionUpload == "" {
		return ClientVersion{}, errors.New("no version found")
	}

	return cv, nil
}

func AllLatestVersions() (ChannelsClientVersions, error) {
	ccvs := make(ChannelsClientVersions, 0)

	for _, c := range Channels {
		cv, err := LatestVersion(c)
		if err != nil {
			return ChannelsClientVersions{}, err
		}

		ccvs[c] = cv
	}

	return ccvs, nil
}
