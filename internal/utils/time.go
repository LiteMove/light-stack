package utils

import (
	"fmt"
	"time"
)

// ParseToTime 将字符串解析为时间指针，支持多种格式
func ParseToTime(timeStr string) (*time.Time, error) {
	if timeStr == "" {
		return nil, nil
	}
	// 尝试解析多种时间格式
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, timeStr)
		if err == nil {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("无法解析时间格式: %s", timeStr)
}
