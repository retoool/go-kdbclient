package kdb

import (
	"fmt"
	"go-datacalc/utils/kdb/entity"
	"net/http"
	"time"
)

// DeteleMetricRange 根据开始时间和结束时间删除指定名称的指标
func DeteleMetricRange(host, port string, pointname string, starttime time.Time, endtime time.Time) *http.Response {
	beginunix := starttime.UnixMilli()
	endUnix := endtime.UnixMilli()
	k := entity.NewKairosdb(host, port)
	bodytext := make(map[string]interface{})

	bodytext = map[string]interface{}{
		"start_absolute": beginunix,
		"end_absolute":   endUnix,
		"metrics": []map[string]interface{}{
			{
				"name": pointname,
			},
		},
	}
	response, err := entity.PostRequest(k.DeleteUrl, bodytext, k.Headersjson)
	if err != nil {
		fmt.Println(err)
	}
	return response
}

// DeteleMetric 根据名称删除指标
func DeteleMetric(host, port string, pointName string) *http.Response {
	k := entity.NewKairosdb(host, port)
	req, err := http.NewRequest("DELETE", k.DelUrl+pointName, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return response
}
