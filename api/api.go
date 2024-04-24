package api

import (
	"net/http"

	"github.com/retoool/go-kdbclient/entity"
)

type KdbClient interface {
	AddCustomFunc(name, funcJsStr string)
	AddPoint(pointname string, tags []string, aggr string, aligntime string, minvalue string,
		maxvalue string, samplingValue string, samplingUnit string)
	AddPoints(pointnames []string, tags []string, aggr string, aligntime string, minvalue string,
		maxvalue string, samplingValue string, samplingUnit string)
	Query() (entity.Response_trans, error)
	Delete() (*http.Response, error)
}

type Response_trans interface {
	GetResult() []map[string]map[int64]float64
	SetResult(result []map[string]map[int64]float64)
	GetCode() int
	SetCode(code int)
	GetMessages() []string
	SetMessages(messages []string)
}
