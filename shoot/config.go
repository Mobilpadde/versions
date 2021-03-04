package shoot

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"versions/execute"

	"github.com/raff/godet"
)

type Shooter struct {
	chrome *exec.Cmd
}

func New(verboser bool) *Shooter {
	chrome := execute.Command(verboser, "google-chrome-beta", "", []string{}, "--remote-debugging-port=9222")
	chrome.Start()

	return &Shooter{
		chrome: chrome,
	}
}

func (s *Shooter) Close() {
	s.chrome.Process.Kill()
}

func Shoot(uri, path, sha1 string, wait int) {
	remote, err := godet.Connect("localhost:9222", false)
	if err != nil {
		log.Println(err)
		return
	}
	defer remote.Close()

	tab, _ := remote.NewTab(uri)
	defer remote.CloseTab(tab)

	time.Sleep(time.Second * 2)
	remote.Reload()
	time.Sleep(time.Second * 2)

	remote.SaveScreenshot(fmt.Sprintf("%s/%s.png", path, sha1), 0644, 100, true)
	time.Sleep(time.Second)
}
