package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Mobilpadde/versions/ani"
	"github.com/Mobilpadde/versions/execute"
	"github.com/Mobilpadde/versions/git"
	"github.com/Mobilpadde/versions/logs"
	"github.com/Mobilpadde/versions/shoot"
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
	var paths string
	var port int
	var wait int
	var commits int
	var verbose bool
	var verboser bool

	flag.StringVar(&repo, "repo", "", "the path of a git-repo")
	flag.StringVar(&dump, "dump", "./screendumps", "the path the screendumps-dir")
	flag.StringVar(&out, "out", "./out", "the path for the generated .webp's")
	flag.StringVar(&manager, "manager", "pnpm", "which node-package manager to use")
	flag.StringVar(&cmd, "cmd", "dev", "the node-package manager (`-manager`) command used to run the dev-server")
	flag.StringVar(&installCmd, "install", "i", "the package manager install command (like `pnpm i`)")
	flag.StringVar(&paths, "paths", "/", "the path to screenshot, comma-seperated")
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

	os.Mkdir(dump, 0777)
	os.Mkdir(out, 0777)

	logsData := logs.GetLogs(repo)

	if commits > 0 {
		logsData = logsData[:commits]
	}

	shooter := shoot.New(verboser)
	defer func() {
		log.Println(git.ChangeCommit(repo, "master"))
		log.Println(git.ChangeCommit(repo, "main"))
	}()

	pathsSplit := strings.Split(paths, ",")
	pathRe := regexp.MustCompile(`[^\w+]`)

	for i, l := range logsData {
		log.Printf("Checking out: [%s]: %s", l.SHA1, l.Title)
		git.ChangeCommit(repo, l.SHA1)

		installs := strings.Split(installCmd, " ")
		d := execute.Command(verbose || verboser, manager, repo, []string{}, installs...)
		d.Run()

		cmds := strings.Split(cmd, " ")
		s := execute.Command(verbose || verboser, manager, repo, []string{"PORT=" + strconv.Itoa(port)}, cmds...)
		s.Start()

		time.Sleep(time.Second * time.Duration(wait))
		for _, path := range pathsSplit {
			if path == "" {
				continue
			}

			dumpPath := dump + "/" + pathRe.ReplaceAllString(path, "")
			os.Mkdir(dumpPath, 0777)
			time.Sleep(time.Second)
			shoot.Shoot(fmt.Sprintf("http://localhost:%d%s", port, path), dumpPath, fmt.Sprintf("%d_%s", i, l.SHA1), wait)
			time.Sleep(time.Second)
		}

		for i := 0; i < 2; i++ {
			time.Sleep(time.Millisecond * 250)

			if s.Process != nil {
				s.Process.Kill()

				d := exec.Command("lsof", "-t", "-i:"+strconv.Itoa(port))
				b, _ := d.Output()
				d.Run()

				pID := string(b)
				if pID != "" {
					exec.Command("kill", pID).Run()
				}
			}
		}

		time.Sleep(time.Millisecond * 250)
	}
	shooter.Close()

	for _, path := range pathsSplit {
		if path == "" {
			continue
		}

		repPath := pathRe.ReplaceAllString(path, "")
		dumpPath := dump + "/" + repPath
		ani.DrawAll(dumpPath, path, logsData)

		if repPath == "" {
			repPath = "index"
		}
		ani.Make(dumpPath, out+"/"+repPath+".webp", logsData)
	}
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
