package entity

type Response struct {
	statusCode int
	Errors     []string `json:"errors,omitempty"`
}

func (r *Response) SetStatusCode(code int) {
	r.statusCode = code
}

func (r *Response) GetStatusCode() int {
	return r.statusCode
}

func (r *Response) GetErrors() []string {
	return r.Errors
}
