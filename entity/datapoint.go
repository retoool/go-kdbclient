package entity

import (
	"encoding/json"
	"errors"
)

// DataPoint 结构体表示一个测量值，存储测量发生的时间和其值。
type DataPoint struct {
	timestamp int64
	value     interface{}
}

// NewDataPoint 创建一个新的 DataPoint 实例。
func NewDataPoint(ts int64, val interface{}) *DataPoint {
	return &DataPoint{
		timestamp: ts,
		value:     val,
	}
}

// Timestamp 返回 DataPoint 实例的时间戳。
func (dp *DataPoint) Timestamp() int64 {
	return dp.timestamp
}

// Int64Value 返回 DataPoint 实例的 int64 类型的值。
func (dp *DataPoint) Int64Value() (int64, error) {
	val, ok := dp.value.(int64)
	if !ok {
		v, ok := dp.value.(int)
		if !ok {
			return 0, ErrorDataPointInt64
		}
		val = int64(v)
	}

	return val, nil
}

// Float64Value 返回 DataPoint 实例的 float64 类型的值。
func (dp *DataPoint) Float64Value() (float64, error) {
	val, ok := dp.value.(float64)
	if !ok {
		return 0, ErrorDataPointFloat64
	}
	return val, nil
}

// Float32Value 返回 DataPoint 实例的 float32 类型的值。
func (dp *DataPoint) Float32Value() (float32, error) {
	val, ok := dp.value.(float32)
	if !ok {
		return 0, ErrorDataPointFloat32
	}
	return val, nil
}

// MarshalJSON 实现了 json.Marshaler 接口，将 DataPoint 实例转换为 JSON 格式。
func (dp *DataPoint) MarshalJSON() ([]byte, error) {
	data := []interface{}{dp.timestamp, dp.value}
	return json.Marshal(data)
}

// UnmarshalJSON 实现了 json.Unmarshaler 接口，将 JSON 格式转换为 DataPoint 实例。
func (dp *DataPoint) UnmarshalJSON(data []byte) error {
	var arr []interface{}
	err := json.Unmarshal(data, &arr)
	if err != nil {
		return err
	}

	var v float64
	ok := false
	if v, ok = arr[0].(float64); !ok {
		return errors.New("Invalid Timestamp type")
	}

	// Update the receiver with the values decoded.
	dp.timestamp = int64(v)
	dp.value = arr[1]

	return nil
}
