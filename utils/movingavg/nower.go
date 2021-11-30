package movingavg

import "time"

type Nower interface {
	Now() time.Time
}

type NowTime struct {
}

func (nowTime NowTime) Now() time.Time {
	return time.Now()
}
