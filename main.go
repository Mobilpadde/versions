package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"versions/execute"
	"versions/gif"
	"versions/git"
	"versions/logs"
	"versions/shoot"
)

func main() {
	var path string
	var port int
	var wait int
	var commits int

	flag.StringVar(&path, "path", "", "the path of a git-repo")
	flag.IntVar(&port, "port", 5000, "port of app")
	flag.IntVar(&wait, "wait", 5, "how long to wait before screendumping")
	flag.IntVar(&commits, "commits", 0, "how many commits to dump")

	flag.Parse()

	if path == "" {
		log.Fatalln(errors.New("please enter a path of a git-repo").Error())
	}

	if err := exists(path); err != nil {
		log.Fatalln(err.Error())
	}

	os.RemoveAll("./screendumps")
	os.RemoveAll("./out.gif")

	os.Mkdir("./screendumps", 0777)

	shooter := shoot.New()
	defer shooter.Close()

	logsData := logs.GetLogs(path)
	if commits > 0 {
		logsData = logsData[:commits]
	}

	defer func(port, commits int, logs []logs.Log) {
		if commits <= 0 {
			port = port + len(logs)
		} else {
			port = port + commits
		}

		dir, _ := os.Getwd()
		for i := range logsData {
			k := execute.Command(false, "make", dir, []string{}, "PORT="+strconv.Itoa(port-i), "-i", "kill")
			k.Run()
			k.Wait()
		}
	}(port, commits, logsData)

	defer git.ChangeCommit(path, "master")
	for i, l := range logsData {
		time.Sleep(time.Second * time.Duration(wait))

		log.Printf("Checking out: %s (%s)", l.Title, string(l.SHA1[:5]))
		git.ChangeCommit(path, l.SHA1)

		d := execute.Command(true, "yarn", path, []string{})
		d.Run()
		d.Wait()

		s := execute.Command(true, "yarn", path, []string{"PORT=" + strconv.Itoa(port)}, "dev")
		s.Start()

		time.Sleep(time.Second * time.Duration(wait))
		shoot.Shoot(fmt.Sprintf("http://localhost:%d", port), "./screendumps", fmt.Sprintf("%d_%s", i, l.SHA1), wait)

		s.Process.Kill()
		port++
	}

	gif.DrawAll("./screendumps", logsData)
	gif.Make("./screendumps", logsData)
}

func exists(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		return err
	}
	return nil
}
