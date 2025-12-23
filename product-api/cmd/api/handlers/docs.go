package handlers

// HTTPError response for api
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

// Message response for api
type MessageWrapper struct {
	Message string `json:"message" example:"OK"`
}
