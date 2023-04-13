package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Kairosdb 结构体
type Kairosdb struct {
	Host        string
	Port        string
	QueryUrl    string
	DeleteUrl   string
	PushUrl     string
	DelUrl      string
	Headersjson map[string]string
	Headersgzip map[string]string
}

// NewKairosdb 创建一个新的Kairosdb实例
func NewKairosdb(host string, port string) *Kairosdb {
	k := Kairosdb{
		Host:        host,
		Port:        port,
		QueryUrl:    fmt.Sprintf("http://%s:%s/api/v1/datapoints/query", host, port),  // 查询URL
		DeleteUrl:   fmt.Sprintf("http://%s:%s/api/v1/datapoints/delete", host, port), // 删除URL
		PushUrl:     fmt.Sprintf("http://%s:%s/api/v1/datapoints", host, port),        // 推送URL
		DelUrl:      fmt.Sprintf("http://%s:%s/api/v1/metric/", host, port),           // 删除指定metric URL
		Headersjson: map[string]string{"content-type": "application/json"},            // JSON请求头
		Headersgzip: map[string]string{"content-type": "application/gzip"},            // GZIP请求头
	}
	return &k
}

// PostRequest 发送POST请求
func PostRequest(url string, bodyText interface{}, headers map[string]string) (*http.Response, error) {
	jsonBody, err := json.Marshal(bodyText) // 将请求体序列化为JSON格式
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody)) // 创建POST请求
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value) // 设置请求头
	}
	client := &http.Client{}
	response, err := client.Do(req) // 发送请求
	if err != nil {
		return nil, err
	}
	return response, nil
}
