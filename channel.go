package main

type Channel int

const (
	LIVE Channel = iota
	ZIntegration
	ZCanary
	ZFlag
	ZNext
)

var Channels = []Channel{LIVE, ZIntegration, ZCanary, ZFlag, ZNext}

func (c Channel) String() string {
	switch c {
	case LIVE:
		return "LIVE"
	case ZIntegration:
		return "ZIntegration"
	case ZCanary:
		return "ZCanary"
	case ZFlag:
		return "ZFlag"
	case ZNext:
		return "ZNext"
	}

	return "unknown"
}
