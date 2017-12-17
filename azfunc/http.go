package azfunc

import "io"

// HTTPRequest contains the HTTP request received from the Azure Functions runtime
type HTTPRequest struct {
	Method     string
	URL        string
	Headers    map[string]string
	Body       io.ReadCloser
	Params     map[string]string
	StatusCode string
	Query      map[string]string
	IsRaw      bool
}
