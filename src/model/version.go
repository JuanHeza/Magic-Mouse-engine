package model

// VersionResponseBody /version response body
type VersionResponseBody struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"buid_date"`
}

// VersionResponse Response holder
type VersionResponse struct {
	Body VersionResponseBody
	Err  error
}

// VersionRequest Request holder
type VersionRequest struct {
}
