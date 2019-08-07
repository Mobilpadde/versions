package git

import (
	"log"

	"versions/execute"
)

func GetLogs(path string) string {
	cmd := execute.Command(false, "git", path, []string{}, "--no-pager", "log", "--pretty=oneline")

	logs, err := cmd.Output()
	if err != nil {
		log.Panicln(err.Error())
	}

	return string(logs)
}

func ChangeCommit(path, sha1 string) string {
	cmd := execute.Command(false, "git", path, []string{}, "checkout", sha1)

	logs, err := cmd.Output()
	if err != nil {
		log.Panicln(err.Error())
	}

	return string(logs)
}
