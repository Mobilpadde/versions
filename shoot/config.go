package shoot

import (
	"fmt"
	"os/exec"
	"time"
	"versions/execute"

	"github.com/raff/godet"
)

type Shooter struct {
	chrome *exec.Cmd
}

func New() *Shooter {
	chrome := execute.Command(false, "chromium", "", []string{}, "--remote-debugging-port=9222")
	chrome.Start()

	return &Shooter{
		chrome: chrome,
	}
}

func (s *Shooter) Close() {
	s.chrome.Process.Kill()
}

func Shoot(uri, path, sha1 string, wait int) {
	remote, _ := godet.Connect("localhost:9222", false)
	defer remote.Close()

	remote.SetVisibleSize(1920, 1750)

	tab, _ := remote.NewTab(uri)
	defer remote.CloseTab(tab)

	time.Sleep(time.Second * time.Duration(wait))
	remote.SaveScreenshot(fmt.Sprintf("%s/%s.png", path, sha1), 0644, 100, true)
	time.Sleep(time.Second)
}
