package retry

import (
	"github.com/smallnest/rpcx/log"
	"time"
)

/**
最多重试attempts次，如果调用fn返回错误，
等待sleep的时间，而下次错误重试就需要等待两倍的时间了。
还有一点是错误的类型，常规错误会重试，而stop类型的错误会中断重试，
这也提供了一种中断机制。
*/
func Retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if s, ok := err.(stop); ok {
			return s.error
		}

		if attempts--; attempts > 0 {
			log.Warnf("retry func error: %s. attemps #%d after %s.", err.Error(), attempts, sleep)
			time.Sleep(sleep)
			return Retry(attempts, 2*sleep, fn)
		}
		return err
	}
	return nil
}

type stop struct {
	error
}

func NoRetryError(err error) stop {
	return stop{err}
}
