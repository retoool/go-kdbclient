package kdb

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/retoool/go-kdbclient/entity"
	"net/http"
	"strconv"
	"strings"
)

func PushMsgToKdb(host, port string, msg []string) *http.Response {
	var datas []*SensorData
	for i := 0; i < len(msg); i++ {
		d, err := ParseSensorData(msg[i])
		if err != nil {
			fmt.Println(err)
			continue
		}
		datas = append(datas, d)
	}
	jsonData, err := json.Marshal(datas)
	if err != nil {
		fmt.Println(err)
	}
	//将 JSON 数据压缩
	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	if _, err := gz.Write(jsonData); err != nil {
		fmt.Println(err)
	}
	if err := gz.Close(); err != nil {
		fmt.Println(err)
	}
	k := entity.NewKairosdb(host, port)
	req, err := http.NewRequest("POST", k.PushUrl, &compressed)
	//req, err := http.NewRequest("POST", k.PushUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/gzip")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	return response
}

type SensorData struct {
	Name      string            `json:"name"`
	Timestamp int               `json:"timestamp"`
	Value     float64           `json:"value"`
	Tags      map[string]string `json:"tags"`
}

func ParseSensorData(str string) (*SensorData, error) {

	parts := strings.Split(str, "@F:")

	index := strings.LastIndex(parts[0], ":")
	devName := parts[0][:index]
	pointName := parts[0][index+1:]
	parts2 := strings.Split(parts[1], ":")
	value := parts2[0]
	time := parts2[1]
	timeint, err := strconv.Atoi(time)
	if err != nil {
		fmt.Println(err)
	}
	valuefloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Println(err)
	}
	data := &SensorData{
		Name:      pointName,
		Timestamp: timeint,
		Value:     valuefloat,
		Tags: map[string]string{
			"project": devName,
		},
	}
	return data, nil
}
