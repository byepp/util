package looputil

import (
	"time"
)

// 执行函数如果失败则重试
// 次数,间隔,执行内容
func Retry(attempts int, sleep time.Duration, fn func() error) (err error) {
	for ; attempts > 0; attempts-- {
		if err = fn(); err == nil {
			// 执行成功返回
			break
		}
		time.Sleep(sleep)
	}
	return err
}

// 每天固定时间执行一次函数
func StartDailyFunc(hour, min int, fn func()) {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location())
			if now.After(next) {// 如果当前已经过了今天这个定时的时间
				// 跳到明天的这个时间
				next = next.Add(time.Hour * 24)
			}
			t := time.NewTimer(next.Sub(now))
			<-t.C
			fn()
			// 等待以防止 fn 执行时间过短导致一天多次执行
			time.Sleep(10 * time.Microsecond)
		}
	}()
}

// 每隔一段时间就执行一次函数
func StartIntervalFunc(interval time.Duration, fn func()) {
	go func() {
		for {
			<-time.NewTimer(interval).C
			fn()
		}
	}()
}
