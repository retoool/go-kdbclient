package kdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/retoool/go-kdbclient/entity"
)

// KdbClient 是一个Kairosdb客户端
type KdbClient struct {
	Bodytext map[string]interface{} // 请求体
	Kdbhttp  *entity.Kairosdb       // Kairosdb实例
}

// NewClient 创建一个新的KdbClient
func NewClient(host string, port string, starttime time.Time, endtime time.Time) *KdbClient {
	var client KdbClient
	beginunix := starttime.UnixMilli()
	endUnix := endtime.UnixMilli()
	client.Bodytext = map[string]interface{}{
		"start_absolute": beginunix,
		"end_absolute":   endUnix,
		"metrics":        []map[string]interface{}{},
	}
	client.Kdbhttp = entity.NewKairosdb(host, port)
	return &client
}

// AddPoint 添加一个数据点
func (client *KdbClient) AddPoint(pointname string, tags []string, aggr string, aligntime string, minvalue string,
	maxvalue string, samplingValue string, samplingUnit string) *KdbClient {
	if samplingValue == "" && samplingUnit == "" {
		samplingValue = "10"
		samplingUnit = "years"
	}
	metric := make(map[string]interface{})
	metric["group_by"] = []map[string]interface{}{
		{"name": "tag", "tags": []string{"project"}},
	}
	metric["name"] = pointname
	metric["tags"] = map[string]interface{}{
		"project": tags,
	}
	aggregators := make([]interface{}, 0)
	if minvalue != "" {
		minvalueInt, err := strconv.Atoi(minvalue)
		if err != nil {
			fmt.Println(err)
		}
		minAggregator := map[string]interface{}{
			"name":      "filter",
			"filter_op": "lt",
			"threshold": minvalueInt,
		}
		aggregators = append(aggregators, minAggregator)
	}
	if maxvalue != "" {
		maxvalueInt, err := strconv.Atoi(maxvalue)
		if err != nil {
			fmt.Println(err)
		}
		maxAggregator := map[string]interface{}{
			"name":      "filter",
			"filter_op": "gt",
			"threshold": maxvalueInt,
		}
		aggregators = append(aggregators, maxAggregator)
	}
	if aggr != "" {
		newAggregator := map[string]interface{}{
			"name":     aggr,
			"sampling": map[string]string{"value": samplingValue, "unit": samplingUnit},
		}
		if aligntime == "start" {
			newAggregator["align_start_time"] = true
		} else if aligntime == "end" {
			newAggregator["align_end_time"] = true
		} else if aligntime == "sample" {
			newAggregator["align_sampling"] = true
		}
		aggregators = append(aggregators, newAggregator)
	}
	metric["aggregators"] = aggregators
	client.Bodytext["metrics"] = append(client.Bodytext["metrics"].([]map[string]interface{}), metric)
	return client
}

// AddPoints 添加多个数据点
func (client *KdbClient) AddPoints(pointnames []string, tags []string, aggr string, aligntime string, minvalue string,
	maxvalue string, samplingValue string, samplingUnit string) {
	for _, pointname := range pointnames {
		if samplingValue == "" && samplingUnit == "" {
			samplingValue = "10"
			samplingUnit = "years"
		}
		metric := make(map[string]interface{})
		metric["group_by"] = []map[string]interface{}{
			{"name": "tag", "tags": []string{"project"}},
		}
		metric["name"] = pointname
		metric["tags"] = map[string]interface{}{
			"project": tags,
		}
		aggregators := make([]interface{}, 0)
		if minvalue != "" {
			minAggregator := map[string]interface{}{
				"name":      "filter",
				"filter_op": "lt",
				"threshold": minvalue,
			}
			aggregators = append(aggregators, minAggregator)
		}
		if maxvalue != "" {
			maxAggregator := map[string]interface{}{
				"name":      "filter",
				"filter_op": "gt",
				"threshold": maxvalue,
			}
			aggregators = append(aggregators, maxAggregator)
		}
		if aggr != "" {
			newAggregator := map[string]interface{}{
				"name":     aggr,
				"sampling": map[string]string{"value": samplingValue, "unit": samplingUnit},
			}
			if aligntime == "start" {
				newAggregator["align_start_time"] = true
			} else if aligntime == "end" {
				newAggregator["align_end_time"] = true
			}
			aggregators = append(aggregators, newAggregator)
		}
		metric["aggregators"] = aggregators
		client.Bodytext["metrics"] = append(client.Bodytext["metrics"].([]map[string]interface{}), metric)
	}

}

// Query 查询数据
func (client *KdbClient) Query() []map[string]map[int64]float64 {
	response, err := entity.PostRequest(client.Kdbhttp.QueryUrl, client.Bodytext, client.Kdbhttp.Headersjson)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close() // 优化：关闭response.Body
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	resp := entity.NewQueryResponse(response.StatusCode)
	err = json.Unmarshal(contents, resp)
	if err != nil {
		fmt.Println(err)
	}
	qrMap := make([]map[string]map[int64]float64, len(resp.QueriesArr))
	if len(resp.QueriesArr) == 0 {
		fmt.Print("kairosdb返回数据异常, ")
		code := resp.GetStatusCode()
		errors := resp.GetErrors()
		fmt.Println(code, errors)
		return nil
	}
	for i := range resp.QueriesArr {
		qrMap[i] = make(map[string]map[int64]float64)
		for j := range resp.QueriesArr[i].ResultsArr {
			results := resp.QueriesArr[i].ResultsArr[j]
			points := results.DataPoints
			if len(results.Tags["project"]) <= 0 {
				fmt.Println(results.Name + ",未查询到数据")
				continue
			}
			tag := results.Tags["project"][0]
			if len(points) == 0 {
				fmt.Println(tag + ":" + results.Name + ",未查询到数据")
				continue
			}
			for y := range points {
				value, err := points[y].Float64Value()
				if err != nil {
					fmt.Println(err)
				}
				scale := math.Pow(10, float64(6))
				value = math.Round(value*scale) / scale

				timestamp := points[y].Timestamp()

				if qrMap[i][tag] == nil {
					qrMap[i][tag] = make(map[int64]float64)
				}
				qrMap[i][tag][timestamp] = value
			}

		}
	}
	return qrMap
}

// Delete 删除数据
func (client *KdbClient) Delete() *http.Response {
	response, err := entity.PostRequest(client.Kdbhttp.DeleteUrl, client.Bodytext, client.Kdbhttp.Headersjson)
	if err != nil {
		fmt.Println(err)
	}
	return response
}
