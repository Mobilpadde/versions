package execute

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Command makes it possible to execute a command
func Command(verbose bool, command, dir string, env []string, vars ...string) *exec.Cmd {
	cmd := exec.Command(command, vars...)
	cmd.Dir = dir
	cmd.Env = append(cmd.Env, env...)

	if verbose {
		outReader, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
			return nil
		}

		outScanner := bufio.NewScanner(outReader)
		go func() {
			for outScanner.Scan() {
				log.Printf("\t > %s\n", outScanner.Text())
			}
		}()

		errReader, err := cmd.StderrPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
			return nil
		}

		errScanner := bufio.NewScanner(errReader)
		go func() {
			for errScanner.Scan() {
				log.Printf("\t > %s\n", errScanner.Text())
			}
		}()
	}

	return cmd
}
