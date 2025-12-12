package metrics

const (
	MetricHTTPRequestsTotal   = "http_requests_total"
	MetricHTTPRequestDuration = "http_request_duration_seconds"
	MetricGRPCRequestsTotal   = "grpc_requests_total"
	MetricGRPCRequestDuration = "grpc_request_duration_seconds"
	MetricGameWSConnections   = "game_ws_connections"
	MetricGameSessionsActive  = "game_sessions_active"
	MetricGameSessionsTotal   = "game_sessions_total"
	MetricGameQuestionsAnswered = "game_questions_answered_total"
	MetricGameButtonPressLatency = "game_button_press_latency_seconds"
)

const (
	HelpHTTPRequestsTotal   = "Total number of HTTP requests"
	HelpHTTPRequestDuration = "HTTP request latency in seconds"
	HelpGRPCRequestsTotal   = "Total number of gRPC requests"
	HelpGRPCRequestDuration = "gRPC request latency in seconds"
	HelpGameWSConnections   = "Number of active WebSocket connections"
	HelpGameSessionsActive  = "Number of currently active game sessions"
	HelpGameSessionsTotal   = "Total number of game sessions created"
	HelpGameQuestionsAnswered = "Total number of questions answered"
	HelpGameButtonPressLatency = "Button press latency from question shown to button pressed"
)

const (
	LabelMethod   = "method"
	LabelEndpoint = "endpoint"
	LabelStatus   = "status"
)

var (
	DefaultHTTPBuckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
	DefaultButtonPressBuckets = []float64{.01, .025, .05, .1, .25, .5, 1, 2.5, 5}
)

const (
	StatusCode2xx = "2xx"
	StatusCode3xx = "3xx"
	StatusCode4xx = "4xx"
	StatusCode5xx = "5xx"
	StatusCodeUnknown = "unknown"
)

