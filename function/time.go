// 时间函数
package function

import (
	"time"

	"github.com/araddon/dateparse"
)

// 获取指定时间所在月的开始 结束时间
func GetMonthStartEnd(t time.Time) (time.Time, time.Time) {
	monthStartDay := t.AddDate(0, 0, -t.Day()+1)
	monthStartTime := time.Date(monthStartDay.Year(), monthStartDay.Month(), monthStartDay.Day(), 0, 0, 0, 0, t.Location())
	monthEndDay := monthStartTime.AddDate(0, 1, -1)
	monthEndTime := time.Date(monthEndDay.Year(), monthEndDay.Month(), monthEndDay.Day(), 23, 59, 59, 0, t.Location())
	return monthStartTime, monthEndTime
}

// 获取指定时间当日开始结束时间
func DayBetweenTime(dateTime time.Time, timeZone string) (beginTime, endTime *time.Time, err error) {
	// 加载时区信息
	local, err := time.LoadLocation(timeZone)
	if err != nil {
		return nil, nil, err
	}
	// 时间设置时区
	dateTime = dateTime.In(local)
	// 日开始结束时间
	beginDate := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 0, 0, 0, 0, dateTime.Location())
	endDate := beginDate.AddDate(0, 0, 1).Add(-time.Second * 1)
	// 成功返回
	return &beginDate, &endDate, nil
}

// 工作日天数增加（去除周六周日）
func WeekDayAdd(dateTime time.Time, number uint32) *time.Time {
	for i := 0; i < int(number); i++ {
		dateTime = dateTime.AddDate(0, 0, 1)
		if dateTime.Weekday() == time.Sunday || dateTime.Weekday() == time.Saturday {
			i--
		}
	}
	return &dateTime
}

// 日期字符串，转时间戳
func Strtotime(dateStr, timeZone string) (*time.Time, error) {
	if location, err := time.LoadLocation(timeZone); err != nil {
		return nil, err
	} else if t, err := dateparse.ParseIn(dateStr, location); err != nil {
		return nil, err
	} else {
		return &t, nil
	}
}
