package pack

import "time"

const (
	SchemeHTTP = "http"
)

const (
	PathPackContent = "/api/packs/%s/content"
	PathPack        = "/api/packs/%s"
)

const (
	DefaultHTTPTimeout = 10 * time.Second
)

