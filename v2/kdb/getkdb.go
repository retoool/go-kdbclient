package kdb

import (
	"encoding/json"
	"fmt"
	"github.com/retoool/go-kdbclient/entity"
	"io/ioutil"
	"strconv"
	"time"
)

// QueryKdb 查询Kairosdb
func QueryKdb(host, port string, pointname string, tags []string, aggr string, starttime time.Time, endtime time.Time,
	aligntime string, minvalue string, maxvalue string, samplingValue string, samplingUnit string) map[string][][]string {
	beginunix := starttime.UnixMilli()
	endUnix := endtime.UnixMilli()
	k := entity.NewKairosdb(host, port)
	if samplingValue == "" && samplingUnit == "" {
		samplingValue = "10"
		samplingUnit = "years"
	}
	bodytext := map[string]interface{}{
		"start_absolute": beginunix,
		"end_absolute":   endUnix,
		"metrics": []map[string]interface{}{
			{
				"group_by": []map[string]interface{}{
					{"name": "tag", "tags": []string{"project"}},
				},
				"name": pointname,
				"tags": map[string]interface{}{
					"project": tags,
				},
				"aggregators": []interface{}{},
			},
		},
	}
	if minvalue != "" {
		minAggregator := map[string]interface{}{
			"name":      "filter",
			"filter_op": "lt",
			"threshold": minvalue,
		}
		bodytext["metrics"].([]map[string]interface{})[0]["aggregators"] = append(
			bodytext["metrics"].([]map[string]interface{})[0]["aggregators"].([]interface{}),
			minAggregator,
		)
	}
	if maxvalue != "" {
		maxAggregator := map[string]interface{}{
			"name":      "filter",
			"filter_op": "gt",
			"threshold": maxvalue,
		}
		bodytext["metrics"].([]map[string]interface{})[0]["aggregators"] = append(
			bodytext["metrics"].([]map[string]interface{})[0]["aggregators"].([]interface{}),
			maxAggregator,
		)
	}
	if aggr != "" {
		newAggregator := map[string]interface{}{
			"name":     aggr,
			"sampling": map[string]string{"value": samplingValue, "unit": samplingUnit},
		}
		bodytext["metrics"].([]map[string]interface{})[0]["aggregators"] = append(
			bodytext["metrics"].([]map[string]interface{})[0]["aggregators"].([]interface{}),
			newAggregator,
		)
	}
	if aligntime == "start" {
		aggregators := bodytext["metrics"].([]map[string]interface{})[0]["aggregators"].([]interface{})
		lastAggregator := aggregators[len(aggregators)-1]
		switch a := lastAggregator.(type) {
		case map[string]interface{}:
			a["align_start_time"] = true
		}
	} else if aligntime == "end" {
		aggregators := bodytext["metrics"].([]map[string]interface{})[0]["aggregators"].([]interface{})
		lastAggregator := aggregators[len(aggregators)-1]
		switch a := lastAggregator.(type) {
		case map[string]interface{}:
			a["align_end_time"] = true
		}
	}
	response, err := entity.PostRequest(k.QueryUrl, bodytext, k.Headersjson)
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	qr := entity.NewQueryResponse(response.StatusCode)
	err = json.Unmarshal(contents, qr)
	if err != nil {
		fmt.Println(err)
	}
	return RespToMap(qr)
}

// RespToMap 将QueryResponse转换为map
func RespToMap(resp *entity.QueryResponse) map[string][][]string {
	qrMap := make(map[string][][]string)
	if len(resp.QueriesArr) == 0 {
		fmt.Print("kairosdb返回数据异常, ")
		code := resp.GetStatusCode()
		errors := resp.GetErrors()
		fmt.Println(code, errors)
		return nil
	}
	for i := 0; i < len(resp.QueriesArr[0].ResultsArr); i++ {
		results := resp.QueriesArr[0].ResultsArr[i]
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
		for y := 0; y < len(points); y++ {
			value, err := points[y].Float64Value()
			valuestr := fmt.Sprintf("%.6f", value)
			if err != nil {
				fmt.Println(err)
			}
			timestamp := points[y].Timestamp()
			timestampstr := strconv.Itoa(int(timestamp))
			qrMap[tag] = append(qrMap[tag], []string{timestampstr, valuestr})
		}
	}
	return qrMap
}
