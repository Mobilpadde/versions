package logs

import (
	"strings"

	"versions/git"
)

// Log contains a single git log entry
type Log struct {
	SHA1  string
	Title string
}

// GetLogs gets logs from path
func GetLogs(path string) []Log {
	logSlice := strings.Split(git.GetLogs(path), "\n")
	logSlice = logSlice[:len(logSlice)-1]
	logs := parseLogs(logSlice)

	return logs
}

func parseLogs(logSlice []string) []Log {
	logs := make([]Log, 0)

	for _, s := range logSlice {
		split := strings.SplitN(s, " ", 2)
		logs = append(logs, Log{
			SHA1:  split[0],
			Title: split[1],
		})
	}

	return logs
}
