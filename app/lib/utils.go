package lib

import (
	"fmt"
	"strings"
	"time"
	"unicode"
)

func SecondsUntilMidnight() int64 {
	now := time.Now()

	// 获取当前日期
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 获取下一天的日期
	tomorrow := today.Add(24 * time.Hour)

	// 计算剩余时间
	duration := tomorrow.Sub(now)
	return int64(duration.Seconds())
}

func GetCurYearAndMonth() string {
	now := time.Now()
	return fmt.Sprintf("%d%d", now.Year(), now.Month())
}

func GetCurDay() int64 {
	now := time.Now()
	return int64(now.Day())
}

func GetCurYearAndMonthAndDay() string {
	now := time.Now()
	return fmt.Sprintf("%d%d%d", now.Year(), now.Month(), now.Day())
}

func TimeHasExp(str string) (bool, error) {
	now := time.Now()
	timeA, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		return false, err
	}

	return timeA.Before(now), nil
}

func GetCurTime() string {
	now := time.Now()
	formatted := now.Format("15:04:05")
	return formatted
}

func GetCurTimeDetail(year int, month int, day int) string {
	t := time.Now().AddDate(year, month, day)
	formatted := t.Format("2006-01-02 15:04:05")
	return formatted
}

func NextSaturday() time.Time {
	now := time.Now() // 获取当前时间
	for {
		// 检查当前日期是否是周六
		if now.Weekday() == time.Saturday {
			return now // 如果是周六，则返回当前日期
		}
		// 如果不是周六，则将时间向前推一天直到找到周六为止
		now = now.AddDate(0, 0, 1)
	}
}

func CalDays(t1, t2 string) (int, error) {
	date1, err := time.Parse("2006-01-02", t1)
	if err != nil {
		return 0, err
	}
	date2, err := time.Parse("2006-01-02", t2)
	if err != nil {
		return 0, err
	}
	return int(date2.Sub(date1).Hours()) / 24, nil
}

func IsAllLetters(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false // 发现非字母字符，返回false
		}
	}
	return true
}

func ProcessingCommands(str string) string {
	if IsAllLetters(str) {
		str = strings.ToUpper(str)
	}
	return str
}

func GetUnix(year int, month int, day int) int64 {
	t := time.Now().AddDate(year, month, day)
	return t.Unix()
}

func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
