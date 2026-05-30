package model

// HealthzResponseBody Response holder
type HealthzResponseBody struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

// HealthzRequest Request holder
type HealthzRequest struct {
}

// HealthzResponse Response holder
type HealthzResponse struct {
	Body HealthzResponseBody
	Err  error
}
