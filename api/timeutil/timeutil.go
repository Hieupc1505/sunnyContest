package timeutil

import (
	"time"
)

func SinceFromToSecond(timestamp time.Time) float64 {
	// Lấy thời gian hiện tại
	currentTime := time.Now()

	// Tính toán sự khác biệt theo giây
	duration := currentTime.Sub(timestamp)
	return duration.Seconds()
}
