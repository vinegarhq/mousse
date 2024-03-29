module github.com/vinegarhq/mousse

go 1.22.0

toolchain go1.22.1

require (
	github.com/apprehensions/rbxbin v0.0.0-20240311165649-265a17b5c75d
	github.com/apprehensions/rbxweb v0.0.0-20240309193157-ac2821c3a715
	github.com/diamondburned/arikawa/v3 v3.3.3
	github.com/robloxapi/rbxdhist v0.6.0
)

require (
	github.com/gorilla/schema v1.2.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/robloxapi/rbxver v0.3.0 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
)

replace github.com/apprehensions/rbxweb => ../rbxweb
replace github.com/apprehensions/rbxbin => ../rbxbin
