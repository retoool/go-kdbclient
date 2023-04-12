package entity

// GroupResult 用于存储分组结果
type GroupResult struct {
	Name string `json:"name,omitempty"` // 名称
}

// Results 用于存储查询结果
type Results struct {
	Name       string              `json:"name,omitempty"`     // 名称
	DataPoints []DataPoint         `json:"values,omitempty"`   // 数据点
	Tags       map[string][]string `json:"tags,omitempty"`     // 标签
	Group      []GroupResult       `json:"group_by,omitempty"` // 分组
}

// Queries 用于存储查询
type Queries struct {
	SampleSize int64     `json:"sample_size,omitempty"` // 样本大小
	ResultsArr []Results `json:"results,omitempty"`     // 结果数组
}

// QueryResponse 用于存储查询响应
type QueryResponse struct {
	*Response
	QueriesArr []Queries `json:"queries",omitempty` // 查询数组
}

// NewQueryResponse 用于创建新的查询响应
func NewQueryResponse(code int) *QueryResponse {
	qr := &QueryResponse{
		Response: &Response{},
	}

	qr.SetStatusCode(code)
	return qr
}
