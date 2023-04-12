package entity

import "errors"

var (
	// Metric Errors.
	ErrorMetricNameInvalid = errors.New("指标名称为空") // 指标名称为空
	ErrorTagNameInvalid    = errors.New("标签名称为空") // 标签名称为空
	ErrorTagValueInvalid   = errors.New("标签值为空")  // 标签值为空
	ErrorTTLInvalid        = errors.New("TTL值无效") // TTL值无效

	// Data Point Errors.
	ErrorDataPointInt64   = errors.New("不是int64数据值")   // 不是int64数据值
	ErrorDataPointFloat32 = errors.New("不是float32数据值") // 不是float32数据值
	ErrorDataPointFloat64 = errors.New("不是float64数据值") // 不是float64数据值

	// Query Metric Errors.
	ErrorQMetricNameInvalid     = errors.New("查询指标名称为空")    // 查询指标名称为空
	ErrorQMetricTagNameInvalid  = errors.New("查询指标标签名称为空")  // 查询指标标签名称为空
	ErrorQMetricTagValueInvalid = errors.New("查询指标标签值为空")   // 查询指标标签值为空
	ErrorQMetricLimitInvalid    = errors.New("查询指标限制必须>=0") // 查询指标限制必须>=0

	// Query Builder Errors.
	ErrorAbsRelativeStartSet      = errors.New("绝对和相对开始时间都不能设置") // 绝对和相对开始时间都不能设置
	ErrorRelativeStartTimeInvalid = errors.New("相对开始时间必须>0")     // 相对开始时间必须>0
	ErrorAbsRelativeEndSet        = errors.New("绝对和相对结束时间都不能设置") // 绝对和相对结束时间都不能设置
	ErrorRelativeEndTimeInvalid   = errors.New("相对结束时间必须>0")     // 相对结束时间必须>0
	ErrorStartTimeNotSpecified    = errors.New("未指定开始时间")        // 未指定开始时间
)
