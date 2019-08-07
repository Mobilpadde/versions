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
	split := strings.Split(git.GetLogs(path), "\n\n")
	split = split[:len(split)-1]
	logs := parseLogs(split)

	return logs
}

func parseLogs(logSlice []string) []Log {
	logs := make([]Log, 0)

	for _, s := range logSlice {
		split := strings.Split(s, "\n")
		logs = append(logs, Log{
			SHA1:  split[0],
			Title: split[1],
		})
	}

	return logs
}
