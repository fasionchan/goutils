package gonumplot

import "errors"

var (
	// errEmptyChart 表示 Chart 不含可绘制的 series。
	errEmptyChart = errors.New("chart has no series to render")
)
