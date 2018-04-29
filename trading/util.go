package trading

import "time"

func TimeToString(t time.Time) string {
	return t.Format("2006-01-02")
}
