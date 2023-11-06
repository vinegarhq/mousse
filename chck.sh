#!/bin/sh

# All known used roblox channels, data taken from Latte softworks
set -- \
	live \
	zavatarrelease \
	zavatarteam \
	zavatarteam2 \
	zbisect2 \
	zcanary \
	zcanary1 \
	zcanary2 \
	zcanary3 \
	zcanary4 \
	zcanary5 \
	zcanary6 \
	zcanary7 \
	zcanary8 \
	zcanary9 \
	zcanaryapps \
	zfeatureluasm \
	zfeaturesoundworks-testing \
	zflag \
	zintegration \
	zintegration1 \
	zintegration2 \
	zlive \
	zlive1 \
	zlive2 \
	zlive3 \
	zloom \
	zmacarm64 \
	znext \
	zprojecteliu \
	zprojectluaproxyparts \
	zprojectmovement-alpha \
	zprojectmovement-beta \
	zprojectqt5159 \
	zprojectuwpua \
	zsocialteam \
	zstudioint1 \
	zstudioint2 \
	zwin64_client_test

ck() {
	url="https://clientsettings.roblox.com/v2/client-version/WindowsPlayer/${1:+/channel/$1}"

	out="$(curl -s "$url")"
	err="${out##*message\":}"
	[ "$err" = "$out" ] || {
		printf "%-26s %s\n" "$1" "${err%%\}*}"
		return
	}

	cver="${out##*version\"\:\"}"
	cver="${cver%%\"*}"
	cvu="${out##*"$cver"\"\,\"clientVersionUpload\"\:\"}"
	cvu="${cvu%%\"*}"
	printf '%-26s %-18s %s\n' "$1" "$cver" "$cvu"
}

for c; do 
	ck "$c"
done
