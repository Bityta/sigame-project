package client

import "time"

const (
	MaxMessageSize = 8192
	PongWait       = 60 * time.Second
	JSONPingPeriod = 5 * time.Second
	WriteWait      = 10 * time.Second
)

