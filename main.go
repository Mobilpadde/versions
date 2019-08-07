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
	var repo string
	var dump string
	var out string
	var cmd string
	var port int
	var wait int
	var commits int

	flag.StringVar(&repo, "repo", "", "the path of a git-repo")
	flag.StringVar(&dump, "dump", "./screendumps", "the path the screendumps-dir")
	flag.StringVar(&out, "out", "./out.gif", "the path of the generated gif")
	flag.StringVar(&cmd, "cmd", "dev", "the yarn-command used to run the dev-server")
	flag.IntVar(&port, "port", 5000, "port of app")
	flag.IntVar(&wait, "wait", 5, "how long to wait before screendumping")
	flag.IntVar(&commits, "commits", 0, "how many commits to dump")

	flag.Parse()

	if repo == "" {
		log.Fatalln(errors.New("please enter a path of a git-repo").Error())
	}

	if err := exists(repo); err != nil {
		log.Fatalln(err.Error())
	}

	os.RemoveAll(dump)
	os.RemoveAll(out)

	os.Mkdir(dump, 0777)

	shooter := shoot.New()
	defer shooter.Close()

	logsData := logs.GetLogs(repo)
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

	defer git.ChangeCommit(repo, "master")
	for i, l := range logsData {
		time.Sleep(time.Second * time.Duration(wait))

		log.Printf("Checking out: [%s]: %s", string(l.SHA1[:5]), l.Title)
		git.ChangeCommit(repo, l.SHA1)

		d := execute.Command(true, "yarn", repo, []string{})
		d.Run()
		d.Wait()

		s := execute.Command(true, "yarn", repo, []string{"PORT=" + strconv.Itoa(port)}, cmd)
		s.Start()

		time.Sleep(time.Second * time.Duration(wait))
		shoot.Shoot(fmt.Sprintf("http://localhost:%d", port), dump, fmt.Sprintf("%d_%s", i, l.SHA1), wait)

		s.Process.Kill()
		port++
	}

	gif.DrawAll(dump, logsData)
	gif.Make(dump, out, logsData)
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
