package utils

import "time"

func ConvertToJakartaTime(t time.Time) time.Time {
	loc := time.FixedZone("Asia/Jakarta", 7*60*60)
	return t.In(loc)
}

func TimeNow() time.Time {

	return ConvertToJakartaTime(time.Now())
}
