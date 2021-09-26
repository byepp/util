package timeutil

import (
	"fmt"
	"strconv"
	"time"
)

func Epoch() int64 {
	return time.Now().Unix()
}

func EpochStr() string {
	return fmt.Sprintf("%d", Epoch())
}

func SecToEpoch(sec int64) time.Time {
	return time.Unix(sec, 0)
}

func StrToEpoch(sec string) (time.Time, error) {
	iSec, err := strconv.ParseInt(sec, 10, 64)
	if err != nil {
		return time.Now(), err
	}
	return time.Unix(iSec, 0), nil
}
