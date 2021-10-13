package handcash

import "time"

const isoFormat = "2006-01-02T15:04:05.999Z07:00"

// currentISOTimestamp generates a timestamp in the same format as Javascript Date().toISOString()
func currentISOTimestamp() string {
	return time.Now().UTC().Format(isoFormat)
}
