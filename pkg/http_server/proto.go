package http_server

type healthStatus string

const (
	healthStatusOK   = "OK"
	healthStatusFail = "FAIL"
)

// HealthCheck response is used by health check request
type HealthCheck struct {
	Status healthStatus
}

// StatusResponse response is used by status request
type StatusResponse struct {
	Version    string `json:"version"`
	VersionAPI string `json:"version_api"`
	Build      string `json:"build"`
	Uptime     string `json:"uptime"`
}
