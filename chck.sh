#!/bin/sh

curl https://raw.githubusercontent.com/bluepilledgreat/Roblox-DeployHistory-Tracker/main/ChannelsAll.json \
	| tr -d '[]",' | sed 1d | while read -r channel; do
	url="https://clientsettings.roblox.com/v2/client-version/WindowsPlayer/channel/$channel"
	printf '%-64s' "${channel}:"

	out="$(curl -s "$url")"
	err="${out##*message\":}"
	[ "$err" = "$out" ] || {
		echo "${err%%\}*}"
		continue
	}

	cver="${out##*version\"\:\"}"
	cver="${cver%%\"*}"
	cvu="${out##*"$cver"\"\,\"clientVersionUpload\"\:\"}"
	cvu="${cvu%%\"*}"
	printf '%-18s %s\n' "$cver" "$cvu"
done
