package execute

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Command makes it possible to execute a command
func Command(redirect bool, command, dir string, env []string, vars ...string) *exec.Cmd {
	cmd := exec.Command(command, vars...)
	cmd.Dir = dir
	cmd.Env = append(cmd.Env, env...)

	if redirect {
		reader, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
			return nil
		}

		scanner := bufio.NewScanner(reader)
		go func() {
			for scanner.Scan() {
				log.Printf("\t > %s\n", scanner.Text())
			}
		}()
	}

	return cmd
}
