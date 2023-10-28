package main

import (
	"github.com/vinegarhq/vinegar/roblox/api"
)

type BinaryType int

const (
	WindowsPlayer BinaryType = iota
	WindowsStudio
	WindowsStudio64
	MacPlayer
	MacStudio
)

func (bt BinaryType) BinaryName() string {
	switch bt {
	case WindowsPlayer:
		return "WindowsPlayer"
	case WindowsStudio:
		return "WindowsStudio"
	case WindowsStudio64:
		return "WindowsStudio64"
	case MacPlayer:
		return "MacPlayer"
	case MacStudio:
		return "MacStudio"
	default:
		return "unknown"
	}
}

func (bt BinaryType) String() string {
	return bt.BinaryName()
}

type Version struct {
	GUID string
	Real string
}

func (bt BinaryType) Version(channel string) (Version, error) {
	var cv api.ClientVersion

	ep := "v2/client-version/" + bt.String() + "/channel/" + channel
	err := api.Request("GET", "clientsettings", ep, &cv)
	if err != nil {
		return Version{}, err
	}

	return Version{
		GUID: cv.ClientVersionUpload,
		Real: cv.Version,
	}, nil
}
