package kdb

import (
	"github.com/retoool/go-kdbclient/entity"
	"net/http"
)

type PushData struct {
	Name       string            `json:"name"`
	DataPoints [][]interface{}   `json:"datapoints"`
	Tags       map[string]string `json:"tags"`
}

func NewPushData(devName string, dataPoints [][]interface{}, tags map[string]string) PushData {
	p := PushData{
		Name:       devName,
		DataPoints: dataPoints,
		Tags:       tags,
	}
	return p
}
func PushMap(host, port string, datas map[string]map[string]map[int64]float64) (*http.Response, error) {
	var bodys []PushData
	for pointName := range datas {
		for devCode := range datas[pointName] {
			var dataPoints [][]interface{}
			for timestamp := range datas[pointName][devCode] {
				value := datas[pointName][devCode][timestamp]
				dataPoint := []interface{}{timestamp, value}
				dataPoints = append(dataPoints, dataPoint)
			}
			tags := make(map[string]string)
			tags["project"] = devCode
			body := NewPushData(pointName, dataPoints, tags)
			bodys = append(bodys, body)
		}
	}
	k := entity.NewKairosdb(host, port)
	response, err := entity.PostRequest(k.PushUrl, bodys, k.Headersjson)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func PushOnePoint(host, port string, pointName string, datas map[string]map[int64]float64) (*http.Response, error) {
	var bodys []PushData
	for devCode := range datas {
		var dataPoints [][]interface{}
		for timestamp := range datas[devCode] {
			value := datas[devCode][timestamp]
			dataPoint := []interface{}{timestamp, value}
			dataPoints = append(dataPoints, dataPoint)
		}
		tags := make(map[string]string)
		tags["project"] = devCode
		body := NewPushData(pointName, dataPoints, tags)
		bodys = append(bodys, body)
	}
	k := entity.NewKairosdb(host, port)
	response, err := entity.PostRequest(k.PushUrl, bodys, k.Headersjson)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func PushOneValue(host, port string, pointName string, devCode string, timestamp int64, value float64) (*http.Response, error) {
	var bodys []PushData
	tags := make(map[string]string)
	tags["project"] = devCode
	var dataPoints [][]interface{}
	dataPoint := []interface{}{timestamp, value}
	dataPoints = append(dataPoints, dataPoint)
	body := NewPushData(pointName, dataPoints, tags)
	bodys = append(bodys, body)
	k := entity.NewKairosdb(host, port)
	response, err := entity.PostRequest(k.PushUrl, bodys, k.Headersjson)
	if err != nil {
		return nil, err
	}
	return response, nil
}
