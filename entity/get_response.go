package entity

// GetResponse 结构体
type GetResponse struct {
	*Response
	Results []string `json:"results,omitempty"` // Results 字段
}

// NewGetResponse 函数
func NewGetResponse(code int) *GetResponse {
	gr := &GetResponse{
		Response: &Response{},
		Results:  nil,
	}
	gr.SetStatusCode(code)
	return gr
}

// GetResults 函数
func (gr *GetResponse) GetResults() []string {
	return gr.Results
}
