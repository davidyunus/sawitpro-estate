package helper

import "time"

var (
	loc *time.Location
	Now = now
)

func now() time.Time {
	return time.Now().In(loc)
}

func InitTime() (err error) {
	loc, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return err
	}
	return nil
}
