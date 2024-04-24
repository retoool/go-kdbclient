package entity

type Response_trans struct {
	result   []map[string]map[int64]float64
	code     int
	messages []string
}

func (r *Response_trans) GetResult() []map[string]map[int64]float64 {
	return r.result
}

func (r *Response_trans) SetResult(result []map[string]map[int64]float64) {
	r.result = result
}

func (r *Response_trans) GetCode() int {
	return r.code
}

func (r *Response_trans) SetCode(code int) {
	r.code = code
}

func (r *Response_trans) GetMessages() []string {
	return r.messages
}

func (r *Response_trans) SetMessages(messages []string) {
	r.messages = messages
}
