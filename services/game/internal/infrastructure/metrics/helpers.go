package metrics

func statusCodeString(status int) string {
	if status >= 200 && status < 300 {
		return StatusCode2xx
	} else if status >= 300 && status < 400 {
		return StatusCode3xx
	} else if status >= 400 && status < 500 {
		return StatusCode4xx
	} else if status >= 500 {
		return StatusCode5xx
	}
	return StatusCodeUnknown
}

