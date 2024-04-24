package kdb

import (
	"encoding/json"
	"errors"
	"io"
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
	maxvalue string, samplingValue string, samplingUnit string) {
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
		if aggr == "diff" {
			newAggregator := map[string]interface{}{
				"name": aggr,
			}
			aggregators = append(aggregators, newAggregator)
		} else {
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
	}
	metric["aggregators"] = aggregators
	client.Bodytext["metrics"] = append(client.Bodytext["metrics"].([]map[string]interface{}), metric)
	return
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
			if aggr == "diff" {
				newAggregator := map[string]interface{}{
					"name": aggr,
				}
				aggregators = append(aggregators, newAggregator)
			} else {
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
		}
		metric["aggregators"] = aggregators
		client.Bodytext["metrics"] = append(client.Bodytext["metrics"].([]map[string]interface{}), metric)
	}

}

// Query 查询数据
func (client *KdbClient) Query() (entity.Response_trans, error) {
	var response_trans entity.Response_trans
	response, err := entity.PostRequest(client.Kdbhttp.QueryUrl, client.Bodytext, client.Kdbhttp.Headersjson)
	if err != nil {
		return entity.Response_trans{}, err
	}
	response_trans.SetCode(response.StatusCode)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response_trans, err
	}
	resp := entity.NewQueryResponse(response.StatusCode)
	err = json.Unmarshal(contents, resp)
	if err != nil {
		return response_trans, err
	}
	qrMap := make([]map[string]map[int64]float64, len(resp.QueriesArr))
	response_trans.SetResult(qrMap)
	if len(resp.QueriesArr) == 0 {
		code := resp.GetStatusCode()
		codeStr := strconv.Itoa(code)
		errs := resp.GetErrors()
		return response_trans, errors.New("kairosdb返回数据异常, " + codeStr + errs[0])
	}
	var messages []string
	for i := range resp.QueriesArr {
		qrMap[i] = make(map[string]map[int64]float64)
		for j := range resp.QueriesArr[i].ResultsArr {
			results := resp.QueriesArr[i].ResultsArr[j]
			points := results.DataPoints
			if len(results.Tags["project"]) <= 0 {
				messages = append(messages, results.Name+",未查询到数据")
				continue
			}
			tag := results.Tags["project"][0]
			if len(points) == 0 {
				messages = append(messages, tag+":"+results.Name+",未查询到数据")
				continue
			}
			for y := range points {
				value, err := points[y].Float64Value()
				if err != nil {
					return response_trans, err
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
	response_trans.SetMessages(messages)
	return response_trans, nil
}

// Delete 删除数据
func (client *KdbClient) Delete() (*http.Response, error) {
	response, err := entity.PostRequest(client.Kdbhttp.DeleteUrl, client.Bodytext, client.Kdbhttp.Headersjson)
	return response, err
}
