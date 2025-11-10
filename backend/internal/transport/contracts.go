package transport

type ErrorsResponse struct {
	Errors map[string]string `json:"errors"`
}
