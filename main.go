package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"versions/execute"
	"versions/gif"
	"versions/git"
	"versions/logs"
	"versions/shoot"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func main() {
	var repo string
	var dump string
	var out string
	var manager string
	var cmd string
	var installCmd string
	var path string
	var port int
	var wait int
	var commits int
	var verbose bool
	var verboser bool

	flag.StringVar(&repo, "repo", "", "the path of a git-repo")
	flag.StringVar(&dump, "dump", "./screendumps", "the path the screendumps-dir")
	flag.StringVar(&out, "out", "./out.gif", "the path of the generated gif")
	flag.StringVar(&manager, "manager", "pnpm", "which node-package manager to use")
	flag.StringVar(&cmd, "cmd", "dev", "the node-package manager (-manager) command used to run the dev-server")
	flag.StringVar(&installCmd, "install", "i", "the package manager install command (like `pnpm i`)")
	flag.StringVar(&path, "path", "/", "the path to screenshot")
	flag.IntVar(&port, "port", 5000, "port of app")
	flag.IntVar(&wait, "wait", 5, "how long to wait before screendumping")
	flag.IntVar(&commits, "commits", 0, "how many commits to dump")

	flag.BoolVar(&verbose, "v", false, "log node-manager-commands")
	flag.BoolVar(&verboser, "vvv", false, "all the logs")

	flag.Parse()

	if repo == "" {
		log.Fatalln(errors.New("please enter a path of a git-repo").Error())
	}

	if err := exists(repo); err != nil {
		log.Fatalln(err.Error())
	}

	// os.RemoveAll(dump)
	// os.RemoveAll(out)

	os.Mkdir(dump, 0777)

	logsData := logs.GetLogs(repo)
	if commits > 0 {
		logsData = logsData[:commits]
	}

	gif.DrawAll(dump, path, logsData)
	gif.Make(dump, out, logsData)
	os.Exit(0)

	shooter := shoot.New(verboser)
	defer shooter.Close()

	defer func(port, commits int, logs []logs.Log) {
		if commits <= 0 {
			port = port + len(logs)
		} else {
			port = port + commits
		}
	}(port, commits, logsData)

	defer git.ChangeCommit(repo, "master")
	for i, l := range logsData {
		log.Printf("Checking out: [%s]: %s", l.SHA1, l.Title)
		git.ChangeCommit(repo, l.SHA1)

		installs := strings.Split(installCmd, " ")
		d := execute.Command(verbose || verboser, manager, repo, []string{}, installs...)
		d.Run()
		d.Wait()

		cmds := strings.Split(cmd, " ")
		s := execute.Command(verbose || verboser, manager, repo, []string{"PORT=" + strconv.Itoa(port)}, cmds...)
		s.Start()

		time.Sleep(time.Second * time.Duration(wait))
		shoot.Shoot(fmt.Sprintf("http://localhost:%d%s", port, path), dump, fmt.Sprintf("%d_%s", i, l.SHA1), wait)

		s.Process.Kill()
		k := execute.Command(verboser, "make", "./", []string{}, "PORT="+strconv.Itoa(port), "kill")
		k.Run()
		k.Wait()

		time.Sleep(time.Second * time.Duration(wait))
	}

	gif.DrawAll(dump, path, logsData)
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
