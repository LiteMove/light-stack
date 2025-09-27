package utils

import (
	"fmt"
	"time"

	"github.com/LiteMove/light-stack/internal/config"
)

func ParseToTime(timeStr string) (*time.Time, error) {
	if timeStr == "" {
		return nil, fmt.Errorf("时间字符串为空")
	}
	config := config.Get()
	timeFormat := config.App.TimeZone
	if timeFormat == "" {
		timeFormat = "Asia/Shanghai"
	}
	timeLoc, loadErr := time.LoadLocation(timeFormat)
	if loadErr != nil {
		return nil, fmt.Errorf("加载时区失败: %w", loadErr)
	}

	// 尝试解析多种时间格式，明确指定时区为UTC
	layouts := []struct {
		layout string
		loc    *time.Location
	}{
		{time.RFC3339, timeLoc}, // RFC3339通常包含时区信息
		{"2006-01-02 15:04:05", timeLoc},
		{"2006-01-02", timeLoc},
	}

	var t time.Time
	var err error
	var allErrors []error

	for _, l := range layouts {
		t, err = time.ParseInLocation(l.layout, timeStr, l.loc)
		if err == nil {
			return &t, nil
		}
		allErrors = append(allErrors, err)
	}

	return nil, fmt.Errorf("无法解析时间格式: %s, 错误: %v", timeStr, allErrors)
}
