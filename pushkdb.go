package kdb

import (
	"fmt"
	"go-datacalc/utils/kdb/entity"
	"net/http"
)

type PushData struct {
	Name       string            `json:"name"`
	DataPoints [][]any           `json:"datapoints"`
	Tags       map[string]string `json:"tags"`
}

func NewData(devName string, dataPoints [][]any, tags map[string]string) PushData {
	p := PushData{
		Name:       devName,
		DataPoints: dataPoints,
		Tags:       tags,
	}
	return p
}
func PushMap(host, port string, datas map[string]map[string]map[int64]float64) *http.Response {
	var bodys []PushData
	for pointName := range datas {
		for devCode := range datas[pointName] {
			var dataPoints [][]any
			for timestamp := range datas[pointName][devCode] {
				value := datas[pointName][devCode][timestamp]
				dataPoint := []any{timestamp, value}
				dataPoints = append(dataPoints, dataPoint)
			}
			tags := make(map[string]string)
			tags["project"] = devCode
			body := NewData(pointName, dataPoints, tags)
			bodys = append(bodys, body)
		}
	}
	k := entity.NewKairosdb(host, port)
	response, err := entity.PostRequest(k.PushUrl, bodys, k.Headersjson)
	if err != nil {
		fmt.Println(err)
	}
	return response
}
func PushOnePoint(host, port string, pointName string, datas map[string]map[int64]float64) *http.Response {
	var bodys []PushData
	for devCode := range datas {
		var dataPoints [][]any
		for timestamp := range datas[devCode] {
			value := datas[devCode][timestamp]
			dataPoint := []any{timestamp, value}
			dataPoints = append(dataPoints, dataPoint)
		}
		tags := make(map[string]string)
		tags["project"] = devCode
		body := NewData(pointName, dataPoints, tags)
		bodys = append(bodys, body)
	}
	k := entity.NewKairosdb(host, port)
	response, err := entity.PostRequest(k.PushUrl, bodys, k.Headersjson)
	if err != nil {
		fmt.Println(err)
	}
	return response
}
func PushOneValue(host, port string, pointName string, devCode string, timestamp int64, value float64) *http.Response {
	var bodys []PushData
	tags := make(map[string]string)
	tags["project"] = devCode
	var dataPoints [][]any
	dataPoint := []any{timestamp, value}
	dataPoints = append(dataPoints, dataPoint)
	body := NewData(pointName, dataPoints, tags)
	bodys = append(bodys, body)
	k := entity.NewKairosdb(host, port)
	response, err := entity.PostRequest(k.PushUrl, bodys, k.Headersjson)
	if err != nil {
		fmt.Println(err)
	}
	return response
}
